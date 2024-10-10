package main

import "fmt"

// 1.定义接口1
// 2.定义请求结构体，该请求需要实现接口1
// 3.定义处理请求的接口2，该接口2方法参数必须是请求结构体
// 4.定义实现接口2的处理结构体，该结构体字段类型必须是接口2类型
// 5.创建责任链，从最后一个处理者开始顺序定义
// 6.提交请求，从第一个处理者开始

type Request interface {
	GetName() string
}

// 请求实现

type ConcreteRequest struct {
	name string
	age  int
}

func (c *ConcreteRequest) GetName() string {
	return c.name
}
func (c *ConcreteRequest) GetAge() int {
	return c.age
}

// 定义处理器接口

type Handler interface {
	Handle(Request)
}

// 处理器实现

type ConcreteHandler1 struct {
	next Handler // 定义下一个处理器
}

func (h *ConcreteHandler1) Handle(r Request) {
	if r.GetName() == "handler1" {
		// 处理逻辑
		return
	}

	// 传递给下一处理器
	if h.next != nil {
		h.next.Handle(r)
	}
}

// 处理器实现

type ConcreteHandler2 struct {
	next Handler // 定义下一个处理器
}

func (h *ConcreteHandler2) Handle(r Request) {
	if r.GetName() == "handler2" {
		// 处理逻辑
		fmt.Println("ok")
		return
	}

	// 传递给下一处理器
	if h.next != nil {
		h.next.Handle(r)
	}
}

func main() {
	// 创建责任链
	handler2 := &ConcreteHandler2{}
	handler1 := &ConcreteHandler1{next: handler2}

	// 提交请求
	request := &ConcreteRequest{"handler2", 31}
	handler1.Handle(request)
}
