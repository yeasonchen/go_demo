package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

type Person struct{}

func (p *Person) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.ContentLength > 0 {
		fmt.Print("dsfaf")
		bs, _ := io.ReadAll(r.Body)
		w.Write(bs)
	}
	fmt.Fprintf(w, "Hello astaxie!") // 这个写入到 w 的是输出到客户端的
}

func main() {
	var ss int32
	var pointer *int32
	ss = -6000
	pointer = &ss

	fmt.Println(ss)
	fmt.Println(pointer)

	err := http.ListenAndServe(":9010", &Person{}) // 设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
