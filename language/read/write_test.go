package main

import (
	"log"
	"os"
	"testing"
)

func Test_Write(t *testing.T) {
	// falg: 没有当前文件就创建、可读可写模式打开、追加模式
	f, err := os.OpenFile("abb.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		_ = f.Close()
	}()

	// 该方法默认是从开头写入数据，配合上面的O_APPEND可以实现追加写入
	_, err = f.Write([]byte("EDF"))
	if err != nil {
		log.Fatalln("写入数据错误", err)
	}

}
