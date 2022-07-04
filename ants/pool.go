package main

import (
	"context"
	"github.com/YeasonChen/go_demo/ants/internal"
	"sync"
	"sync/atomic"
	"time"
)

// Pool accepts the tasks from client, it limits the total of goroutines to a given number by recycling goroutines.
type Pool struct {
	// capacity of the pool, a negative value means that the capacity of pool is limitless, an infinite pool is used to
	// avoid potential issue of endless blocking caused by nested usage of a pool: submitting a task to pool
	// which submits a new task to the same pool.
	capacity int32

	// running is the number of the currently running goroutines.
	running int32

	// waiting is the number of goroutines already been blocked on pool.Submit(), protected by pool.lock
	waiting int32

	// lock for protecting the worker queue.
	lock sync.Locker

	// workers is a slice that store the available workers.
	workers workerArray

	// state is used to notice the pool to closed itself.
	state int32

	// cond for waiting to get an idle worker.
	cond *sync.Cond

	// workerCache speeds up the obtainment of a usable worker in function:retrieveWorker.
	workerCache sync.Pool

	// 定期清理过期的worker
	heartbeatDone int32
	stopHeartbeat context.CancelFunc

	options *Options
}

// NewPool generates an instance of ants pool.
func NewPool(size int, options ...Option) (*Pool, error) {
	opts := loadOptions(options...)

	if size <= 0 {
		size = -1
	}

	if expiry := opts.ExpiryDuration; expiry < 0 {
		return nil, ErrInvalidPoolExpiry
	} else if expiry == 0 {
		opts.ExpiryDuration = DefaultCleanIntervalTime
	}

	if opts.Logger == nil {
		opts.Logger = defaultLogger
	}

	p := &Pool{
		capacity: int32(size),
		lock:     internal.NewSpinLock(),
		options:  opts,
	}
	p.workerCache.New = func() interface{} {
		return &goWorker{
			pool: p,
			task: make(chan func(), workerChanCap),
		}
	}
	if p.options.PreAlloc {
		if size == -1 {
			return nil, ErrInvalidPreAllocSize
		}
		p.workers = newWorkerArray(loopQueueType, size)
	} else {
		p.workers = newWorkerArray(stackType, 0)
	}

	p.cond = sync.NewCond(p.lock)

	// Start a goroutine to clean up expired workers periodically.
	var ctx context.Context
	ctx, p.stopHeartbeat = context.WithCancel(context.Background())
	go p.purgePeriodically(ctx)

	return p, nil
}

// ---------------------------------------------------------------------------

// Submit 提交任务到pool中
func (p *Pool) Submit(task func()) error {
	if p.IsClosed() {
		return ErrPoolClosed
	}
	var w *goWorker
	if w = p.retrieveWorker(); w == nil {
		return ErrPoolOverload
	}
	w.task <- task
	return nil
}

// retrieveWorker 获取可用的worker来执行任务，该方法有可能阻塞
func (p *Pool) retrieveWorker() (w *goWorker) {
	// 尝试从缓存中获取worker
	spawnWorker := func() {
		w = p.workerCache.Get().(*goWorker)
		w.run()
	}

	p.lock.Lock()
	w = p.workers.detach()
	if w != nil {
		// 首先从队列中获取worker
		p.lock.Unlock()
	} else if capacity := p.Cap(); capacity == -1 || capacity > p.Running() {
		// 如果工作队列中已经空了，但是还没达到最大的容量，则开启新的worker
		p.lock.Unlock()
		spawnWorker()
	} else {
		// 非阻塞协程池直接返回
		if p.options.Nonblocking {
			p.lock.Unlock()
			return
		}
	retry:
		// 已经阻塞的任务大于外部设置的最大等待数
		if p.options.MaxBlockingTasks != 0 && p.Waiting() >= p.options.MaxBlockingTasks {
			p.lock.Unlock()
			return
		}
		// 走到这里说明需要阻塞这个添加的任务, 等待cond的唤醒
		p.addWaiting(1)
		p.cond.Wait()
		p.addWaiting(-1)

		// 这里表示已经被唤醒
		// 被唤醒的情况：
		// 1、关闭pool的时候
		// 2、worker执行完一个任务之后挪到空闲队列前，发现容量比正在运行的worker数量更小(即容量被调小了)，需要唤醒判断是否要减少阻塞的任务数量
		// 3、worker执行完一个任务挪到空闲队列之后
		// 4、容量被改大，需要新增worker来执行
		// 5、清理垃圾的时候 ？？ 待理解
		// 6、此前的worker发生panic退出

		if p.IsClosed() { // 关闭pool的时候触发 - 情况1
			p.lock.Unlock()
			return
		}

		// 被唤醒但是没有正在执行的worker - 情况5
		var nw int
		if nw = p.Running(); nw == 0 {
			p.lock.Unlock()
			spawnWorker()
			return
		}

		if w = p.workers.detach(); w == nil {
			// 容量被调大的情况会触发 - 情况4
			// 此前的worker发生panic - 情况6
			if nw < p.Cap() {
				p.lock.Unlock()
				spawnWorker()
				return
			}
			// 容量被调小后，容量比正在运行的worker数量更小触发 - 情况2
			goto retry
		}
		// 有新的空闲worker放到队列然后唤醒的情况会直接返回worker - 情况3
		p.lock.Unlock()
	}
	return
}

// Release closes this pool and releases the worker queue.
func (p *Pool) Release() {
	if !atomic.CompareAndSwapInt32(&p.state, OPENED, CLOSED) {
		return
	}
	p.lock.Lock()
	p.workers.reset()
	p.lock.Unlock()
	// There might be some callers waiting in retrieveWorker(), so we need to wake them up to prevent
	// those callers blocking infinitely.
	p.cond.Broadcast()
}

// revertWorker 将空闲的worker放回到pool中
func (p *Pool) revertWorker(worker *goWorker) bool {
	// 工作worker大于设置的容量或者已经关闭pool时通知所有等待的任务并直接返回
	if capacity := p.Cap(); (capacity > 0 && p.Running() > capacity) || p.IsClosed() {
		p.cond.Broadcast()
		return false
	}
	worker.recycleTime = time.Now()
	p.lock.Lock()

	// To avoid memory leaks, add a double check in the lock scope.
	// Issue: https://github.com/panjf2000/ants/issues/113
	if p.IsClosed() {
		p.lock.Unlock()
		return false
	}

	err := p.workers.insert(worker)
	if err != nil {
		p.lock.Unlock()
		return false
	}

	// Notify the invoker stuck in 'retrieveWorker()' of there is an available worker in the worker queue.
	p.cond.Signal()
	p.lock.Unlock()
	return true
}

// ---------------------------------------------------------------------------

func (p *Pool) addRunning(delta int) {
	atomic.AddInt32(&p.running, int32(delta))
}

func (p *Pool) addWaiting(delta int) {
	atomic.AddInt32(&p.waiting, int32(delta))
}

// Running returns the number of workers currently running.
func (p *Pool) Running() int {
	return int(atomic.LoadInt32(&p.running))
}

// Waiting returns the number of tasks which are waiting be executed.
func (p *Pool) Waiting() int {
	return int(atomic.LoadInt32(&p.waiting))
}

// Tune 改变pool的容量，pre-allocation和无限制容量的pool不会受影响
func (p *Pool) Tune(size int) {
	capacity := p.Cap()
	if capacity == -1 || size <= 0 || size == capacity || p.options.PreAlloc {
		return
	}
	atomic.StoreInt32(&p.capacity, int32(size))
	// 如果是改大容量，则需要唤醒一个或者多个等待执行的任务
	// 改小无需额外操作，会在任务执行完成的时候自动缩减worker数
	if size > capacity {
		if size-capacity == 1 {
			p.cond.Signal()
			return
		}
		p.cond.Broadcast()
	}
}

// Cap returns the capacity of this pool.
func (p *Pool) Cap() int {
	return int(atomic.LoadInt32(&p.capacity))
}

// IsClosed indicates whether the pool is closed.
func (p *Pool) IsClosed() bool {
	return atomic.LoadInt32(&p.state) == CLOSED
}

// Free 获取可用的routine数量，-1表示无限制
func (p *Pool) Free() int {
	c := p.Cap()
	if c < 0 {
		return -1
	}
	return c - p.Running()
}

// purgePeriodically clears expired workers periodically which runs in an individual goroutine, as a scavenger.
func (p *Pool) purgePeriodically(ctx context.Context) {
	heartbeat := time.NewTicker(p.options.ExpiryDuration)
	defer func() {
		heartbeat.Stop()
		atomic.StoreInt32(&p.heartbeatDone, 1)
	}()

	for {
		select {
		case <-heartbeat.C:
		case <-ctx.Done():
			return
		}

		if p.IsClosed() {
			break
		}

		// 获取已经超时的worker
		p.lock.Lock()
		expiredWorkers := p.workers.retrieveExpiry(p.options.ExpiryDuration)
		p.lock.Unlock()

		// Notify obsolete workers to stop.
		// This notification must be outside the p.lock, since w.task
		// may be blocking and may consume a lot of time if many workers
		// are located on non-local CPUs.
		for i := range expiredWorkers {
			expiredWorkers[i].task <- nil
			expiredWorkers[i] = nil
		}

		// There might be a situation where all workers have been cleaned up(no worker is running),
		// or another case where the pool capacity has been Tuned up,
		// while some invokers still get stuck in "p.cond.Wait()",
		// then it ought to wake all those invokers.
		if p.Running() == 0 || (p.Waiting() > 0 && p.Free() > 0) {
			p.cond.Broadcast()
		}
	}
}
