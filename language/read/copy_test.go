package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func Test_copy(t *testing.T) {
	f1, err := os.OpenFile("abb.txt", os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		_ = f1.Close()
	}()

	// falg: 没有当前文件就创建、可读可写模式打开
	f2, err := os.OpenFile("abb_copy_1.txt", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		_ = f2.Close()
	}()

	// 使用拷贝方法
	fmt.Println(io.Copy(f2, f1))

	// ioutil方法，不适合拷贝较大的文件
	bs, _ := ioutil.ReadFile("abb.txt")
	ioutil.WriteFile("abb_copy_2.txt", bs, os.ModePerm)

	// 自己实现拷贝方法
	// for {
	// 	total := 0
	// 	buf := make([]byte, 1024)
	// 	n, err := f1.Read(buf)
	// 	if err != nil {
	// 		if err == io.EOF {
	// 			fmt.Println("复制完毕")
	// 			break
	// 		}
	// 		fmt.Println("复制出错，复制了", total)
	// 	}
	// 	total += n
	// 	f2.Write(buf[:n])
	// }
}
