package main

import (
	"fmt"
	"sync"
)

type Singleton struct {
	data string
}

var instance *Singleton
var once sync.Once
var mutex sync.Mutex

func GetInstance() *Singleton {
	mutex.Lock()
	defer mutex.Unlock()

	if instance == nil {
		once.Do(func() {
			instance = &Singleton{data: "Hello, Singleton!"}
		})
	}

	return instance
}

func (s *Singleton) GetData() string {
	return s.data
}

func main() {
	singleton := GetInstance()
	fmt.Println(singleton.GetData()) // 输出: Hello, Singleton!

	// 尝试创建多个实例
	go func() {
		singleton := GetInstance()
		fmt.Println(singleton.GetData()) // 输出: Hello, Singleton!
	}()

	go func() {
		fmt.Println(111)
		singleton := GetInstance()
		fmt.Println(singleton.GetData()) // 输出: Hello, Singleton!
	}()

	// 等待所有goroutine执行完毕
	//mutex.Lock()
	//mutex.Unlock()
}
