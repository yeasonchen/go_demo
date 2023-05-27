package service

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"time"
)

type Option func(*App)

type ShutdownCallback func(ctx context.Context)

func WithShutdownCallbacks(cbs ...ShutdownCallback) Option {
	return func(app *App) {
		app.cbs = cbs
	}
}

type App struct {
	servers []*Server

	// 整个优雅退出的超时时间，默认30s
	shutdownTimeOut time.Duration
	// 优雅退出时候等待处理已有请求时间，默认10s
	waitTime time.Duration
	// 自定义回调超时时间，默认3s
	cbTimeout time.Duration

	// 自定义的回调方法
	cbs []ShutdownCallback
}

func NewApp(servers []*Server, opts ...Option) *App {
	res := &App{
		servers:         servers,
		shutdownTimeOut: 30 * time.Second,
		waitTime:        10 * time.Second,
		cbTimeout:       3 * time.Second,
	}
	for _, opt := range opts {
		opt(res)
	}
	return res
}

// StartAndServe 你主要要实现这个方法
func (app *App) StartAndServe() {
	for _, s := range app.servers {
		srv := s
		go func() {
			if err := srv.Start(); err != nil {
				if err == http.ErrServerClosed {
					log.Printf("服务器%s已关闭", srv.name)
				} else {
					log.Printf("服务器%s异常退出", srv.name)
				}
			}
		}()
	}
	// 从这里开始开始启动监听系统信号
	// ch := make(...) 首先创建一个接收系统信号的 channel ch
	// 定义要监听的目标信号 signals []os.Signal
	// 调用 signal
	ch := make(chan os.Signal, 2)
	signal.Notify(ch, sigs...)
	<-ch

	// 到这里说明收到关闭的信号
	go func() {
		select {
		case <-ch:
			// 再次收到关闭的信号，直接退出
			log.Printf("强制退出")
			os.Exit(1)
		case <-time.After(app.shutdownTimeOut):
			log.Printf("超时强制退出")
			os.Exit(1)
		}
	}()

	app.shutdown()
}

// shutdown APP的关闭流程
func (app *App) shutdown() {
	log.Println("开始关闭应用，停止接收新请求")

	var wg sync.WaitGroup
	wg.Add(len(app.servers))
	for _, s := range app.servers {
		svr := s
		s.rejectReq()
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), app.waitTime)
			if err := svr.stop(ctx); err != nil {
				log.Printf("关闭服务失败%s \n", svr.name)
			}
			cancel()
			wg.Done()
		}()
	}
	wg.Wait()

	log.Println("开始执行自定义回调")
	// 执行回调
	wg.Add(len(app.cbs))
	for _, cb := range app.cbs {
		c := cb
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), app.cbTimeout)
			c(ctx)
			cancel()
			wg.Done()
		}()
	}
	wg.Wait()
	// 释放资源
	log.Println("开始释放资源")
	app.close()
}

func (app *App) close() {
	// 在这里释放掉一些可能的资源
	time.Sleep(time.Second)
	log.Println("应用关闭")
}

type Server struct {
	srv  *http.Server
	name string
	mux  *serverMux

	reqCount  int64
	closeChan chan struct{}
}

type serverMux struct {
	reject bool
	server *Server
	*http.ServeMux
}

func (s *serverMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	atomic.AddInt64(&s.server.reqCount, 1)
	if s.reject {
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte("服务已关闭"))
		return
	}
	s.ServeMux.ServeHTTP(w, r)
	atomic.AddInt64(&s.server.reqCount, -1)
	// 已经处于关闭阶段，需要通知请求已经处理完成
	if atomic.LoadInt64(&s.server.reqCount) == 0 && s.reject {
		close(s.server.closeChan)
	}
}

func NewServer(name string, addr string) *Server {
	mux := &serverMux{ServeMux: http.NewServeMux()}
	ret := &Server{
		name:      name,
		mux:       mux,
		srv:       &http.Server{Addr: addr, Handler: mux},
		closeChan: make(chan struct{}, 1),
	}
	mux.server = ret
	return ret
}

func (s *Server) Handle(pattern string, handler http.Handler) {
	s.mux.Handle(pattern, handler)
}

func (s *Server) Start() error {
	return s.srv.ListenAndServe()
}

func (s *Server) rejectReq() {
	s.mux.reject = true
}

func (s *Server) stop(ctx context.Context) error {
	if atomic.LoadInt64(&s.reqCount) == 0 {
		log.Printf("服务器%s没有正在处理的请求", s.name)
		return s.srv.Shutdown(ctx)
	}
	select {
	case <-s.closeChan:
		log.Printf("服务器%s求处理完成了", s.name)
	case <-ctx.Done():
		log.Printf("服务器%s等待处理完成超时", s.name)
	}
	return s.srv.Shutdown(ctx)
}
