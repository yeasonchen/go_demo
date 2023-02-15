package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	// io.Reader是一个接口 -> Read(p []byte) (n int, err error)
	// file是实现了Reader接口的一个实现
	f, err := os.OpenFile("a.txt", os.O_RDWR, os.ModePerm)
	defer func() {
		_ = f.Close()
	}()
	if err != nil {
		fmt.Println(err)
		return
	}
	// 这里创建一个切片去读取数据，每次最多读取切片的长度的数据直到EOF
	// 如果切片长度是0，则会死循环不停的读取，每次都读取0长度的数据
	// 通常这个长度是1024
	b := make([]byte, 3)
	for {
		n, err := f.Read(b)
		if err != nil {
			if err == io.EOF {
				// 该方法从指定的位置开始读，但是不会改变file内部的游标
				// 后序Read时还是从原有位置开始
				f.ReadAt(b, 1)
				fmt.Println("读完了强行再读1次", string(b))

				n, err := f.Read(b)
				fmt.Println("读完了强行再读2次", string(b), n, err)

				break
			} else {
				fmt.Println("读出错", err)
			}
		}
		// 根据读出的长度打印，因为最后一次读取可能数量不足长度了
		fmt.Printf("读到了数据 %d - %v \n", n, string(b[:n]))
	}

	// 上面自己读数据很麻烦，就有了这个便捷的方法一次性读取，返回切片
	// 即使没有数据也不会返回EOF
	// r, err := ioutil.ReadAll(f)
	// if err != nil {
	// 	if err == io.EOF {
	// 		fmt.Println("读完了")
	// 	} else {
	// 		fmt.Println("读出错", err)
	// 	}
	// }
	// fmt.Printf("读到了数据 - %v \n", string(r))
}
