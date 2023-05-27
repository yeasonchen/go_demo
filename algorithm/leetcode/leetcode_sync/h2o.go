package leetcode_sync

import (
	"fmt"
	"time"
)

// 1117. https://leetcode.cn/problems/building-h2o/
// 输入：HHOOHHHHHHOO
// 输出：HHO  OHH  HHO  HHO
// chan 附带通信功能，里面的数据有含义

var chh = make(chan int)
var cho = make(chan int)

func hydrogen() {
	data := <-chh
	fmt.Printf("H")
	if data == 1 {
		chh <- 2
	} else {
		cho <- 1
	}

}

func oxygen() {
	<-cho
	fmt.Printf("O")
	chh <- 1
}

func H2OExample() {
	var s string
	_, _ = fmt.Scanln(&s)
	for i := 0; i < len(s); i++ {
		if i&3 == 0 {
			go oxygen()
		} else {
			go hydrogen()
		}
	}
	chh <- 1
	time.Sleep(2 * time.Second)
}
