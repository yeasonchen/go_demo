package leetcode_sync

import (
	"fmt"
	"time"
)

// 1114. https://leetcode.cn/problems/print-in-order/

// 多个携程中的3个函数保证顺序输出
// 即下例子中输出应为：
// first...  second...  third...
// first...  second...  third...
// first...  second...  third...

// 思路 -> 定义3个chan，每个函数都依赖一个chan，并且会写入一个chan用于触发下一个函数

const OLEN = 3

func first(ch1, ch2 chan struct{}) {
	<-ch1
	fmt.Println("first...")
	ch2 <- struct{}{}
}

func second(ch1, ch2 chan struct{}) {
	<-ch1
	fmt.Println("second...")
	ch2 <- struct{}{}
}

func third(ch1, ch2 chan struct{}) {
	<-ch1
	fmt.Println("third...")
	ch2 <- struct{}{}
}

func OrderExample() {
	//f, _ := os.Create("./language/trace.out")
	//_ = trace.Start(f)
	//defer trace.Stop()

	ch1 := make(chan struct{}, 1)
	ch2 := make(chan struct{}, 1)
	ch3 := make(chan struct{}, 1)

	for i := 0; i < OLEN; i++ {
		go first(ch3, ch1)
		go second(ch1, ch2)
		go third(ch2, ch3)
	}
	// 触发first执行
	ch3 <- struct{}{}
	time.Sleep(time.Second * 1)
}
