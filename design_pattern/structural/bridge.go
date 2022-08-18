package structural

import "fmt"

// 桥接模式 - 就是要使用接口来抽象，不要依赖具体的实现类
// 现在有发送sms和email的两种方法，并且有发送系统A和发送系统B，每个发送系统中都要能够发送sms和发送email
// 则在先抽象发送SendMessage接口，然后在发送系统中包含接口而不是具体的实现类
// 否则新增发送方法或者新增发送系统都要做很大的修改

// 两种发送消息的方法

type SendMessage interface {
	send(text, to string)
}

type sms struct{}

func NewSms() SendMessage {
	return &sms{}
}

func (*sms) send(text, to string) {
	fmt.Println(fmt.Sprintf("send %s to %s sms", text, to))
}

type email struct{}

func NewEmail() SendMessage {
	return &email{}
}

func (*email) send(text, to string) {
	fmt.Println(fmt.Sprintf("send %s to %s email", text, to))
}

// 两种发送系统

type systemA struct {
	method SendMessage
}

func NewSystemA(method SendMessage) *systemA {
	return &systemA{
		method: method,
	}
}

func (m *systemA) SendMessage(text, to string) {
	m.method.send(fmt.Sprintf("[System A] %s", text), to)
}

type systemB struct {
	method SendMessage
}

func NewSystemB(method SendMessage) *systemB {
	return &systemB{
		method: method,
	}
}

func (m *systemB) SendMessage(text, to string) {
	m.method.send(fmt.Sprintf("[System B] %s", text), to)
}
