package main

import (
	"fmt"
	"log"
)

// 定义一个接口
type Greeter interface {
	Greet(name string) string
}

// 定义一个实现接口的结构体
type EnglishGreeter struct{}

func (g *EnglishGreeter) Greet(name string) string {
	return fmt.Sprintf("Hello, %s!", name)
}

// 定义一个使用接口的结构体
type MyService struct {
	greeter Greeter
}

func (s *MyService) DoSomething(name string) {
	log.Println(s.greeter.Greet(name))
}

func main() {
	// 创建一个 EnglishGreeter 实例
	greeter := &EnglishGreeter{}

	// 创建一个 MyService 实例，并将 greeter 依赖注入
	service := &MyService{greeter}

	// 使用 MyService 实例
	service.DoSomething("John")
}
