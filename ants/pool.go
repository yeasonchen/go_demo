package main

import (
	"errors"
	"sync"
	"time"
)

type sig struct{}

type f func() error

// Pool accept the tasks from client,it limits the total
// of goroutines to a given number by recycling goroutines.
type Pool struct {
	// capacity of the pool.
	capacity int32

	// running is the number of the currently running goroutines.
	running int32

	// expiryDuration set the expired time (second) of every worker.
	expiryDuration time.Duration

	// workers is a slice that store the available workers.
	workers []*Worker

	// release is used to notice the pool to closed itself.
	release chan sig

	// lock for synchronous operation.
	lock sync.Mutex

	once sync.Once
}

func NewPool(size int) (*Pool, error) {
	return NewTimingPool(size, DefaultCleanIntervalTime)
}

func NewTimingPool(size int, expiry time.Duration) (*Pool, error) {
	if size <= 0 {
		return nil, nil
	}
	if expiry <= 0 {
		return nil, nil
	}

	p := &Pool{
		capacity:       int32(size),
		expiryDuration: expiry,
		release:        make(chan sig, 1),
		lock:           sync.Mutex{},
		once:           sync.Once{},
	}

	// 启动定期清理过期worker任务，独立goroutine运行，进一步节省系统资源
	p.monitorAndClear()

	return p, nil
}

func (p *Pool) Submit(task f) error {
	if len(p.release) > 0 {
		return errors.New("协程池已经关闭")
	}
	w := p.getWorker()

	return nil
}

