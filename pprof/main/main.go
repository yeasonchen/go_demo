package main

import (
	"log"
	"os"
	"runtime"
	"runtime/trace"
	"sync"
)

// pprof

//var (
//	cpu = "pprof/main/cpu.out"
//	mem = "pprof/main/mem.out"
//)
//
//func main() {
//
//	// 采样 CPU 运行状态
//	f, err := os.Create(cpu)
//	if err != nil {
//		log.Fatal(err)
//	}
//	_ = pprof.StartCPUProfile(f)
//	defer pprof.StopCPUProfile()
//
//	var wg sync.WaitGroup
//	wg.Add(100)
//	for i := 0; i < 100; i++ {
//		go workOnce(&wg)
//	}
//	wg.Wait()
//
//	// 采样内存状态
//	create, err2 := os.Create(mem)
//	if err2 != nil {
//		log.Fatal(err2)
//	}
//	_ = pprof.WriteHeapProfile(create)
//	_ = create.Close()
//}
//
//func counter() {
//	slice := make([]int, 0)
//	var c int
//	for i := 0; i < 100000; i++ {
//		c = i + 1 + 2 + 3 + 4 + 5
//		slice = append(slice, c)
//	}
//	_ = slice
//}
//
//func workOnce(wg *sync.WaitGroup) {
//	counter()
//	wg.Done()
//}

// trace

var hello []int

func counter(wg *sync.WaitGroup) {
	defer wg.Done()

	slice := []int{0}
	c := 1
	for i := 0; i < 100000; i++ {
		c = i + 1 + 2 + 3 + 4 + 5
		slice = append(slice, c)
	}

	// 防止map被优化
	hello = slice
}

func main() {
	runtime.GOMAXPROCS(1)

	f, err := os.Create("trace.pprof")
	if err != nil {
		log.Fatal(err)
	}
	trace.Start(f)
	defer f.Close()
	defer trace.Stop()

	var wg sync.WaitGroup
	wg.Add(3)
	for i := 0; i < 3; i++ {
		go counter(&wg)
	}
	wg.Wait()
}
