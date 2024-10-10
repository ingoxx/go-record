package main

import "fmt"

//为何指针接收器方法无法为结构体值实现接口

// 方法接收器的类型选择: 值无法调用指针方法，指针方法集不可以让值实现接口，存储在接口的值是拿不到地址
// golang对指针可以隐式的解引用，
type Adult struct {
	Name string
}

//接收者是*Adult，指针类型
//不允许取地址：常量，临时变量，map索引
// 要看对应的值类型是否实现了对应的方法

func (a *Adult) GetName() {
	fmt.Println(a.Name)
}

type Person interface {
	GetName()
}

func main() {
	// var a Person = Adult{"lxb"} 无法传递值类型的结构体值，a的Person接口类型无法寻址Adult{"lxb"}
	var a Person = &Adult{"lxb"}
	a.GetName()
}
