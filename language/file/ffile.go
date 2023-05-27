package main

import (
	"fmt"
	"os"
)

func main() {
	// 文件创建和删除
	// os.Mkdir()		创建文件夹
	// os.MkdirAll()    创建文件夹，以及其子文件夹
	// os.Create()      创建文件
	// os.Remove()      删除文件或文件夹
	// os.RemoveAll()   删除文件夹以及其子文件夹

	// 文件打开后需要关闭
	// Open传入一个filename，返回的f为只读
	f, err := os.Open("aaa.txt")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(f)
	f.Close()

	// 参数：1.文件路径  2.文件打开后的权限  3.文件不存在时创建文件，需要指定权限 ModePerm 为777（1+2+4 = r+w+x）
	f1, err := os.OpenFile("aaa.txt", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(f1)
	f1.Close()

}
