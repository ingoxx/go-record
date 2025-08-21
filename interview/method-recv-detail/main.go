package main

import "fmt"

//方法接收器的值类型与指针类型的细节，问题

type Kind interface {
	Fly()
}

type Brid struct {
	Name string
}

func (b *Brid) Fly() {
	fmt.Println(b.Name)
}

func main() {
	//如果绑定的方法是值类型，以下两种方式均可以赋值给Kind接口
	// var s Kind = Brid{"lxb"}
	// var s1 Kind = &Brid{"lxb"} // golang支持对指针隐式的解引用
	// fmt.Println(s, s1)         //output: {lxb} &{lxb}

	//但如果绑定的方法是指针类型只能用如下方法
	// var s Kind = Brid{"lxb"} 直接报错，在指针类型的方法里，存储在接口的值是拿不到地址的，将值赋值给接口因为接口无法寻址，因此无法实现接口，直接报错
	var s Kind = &Brid{"lxb"}
	fmt.Println(s) //output: &{lxb}

	//那为什么下面这个例子的值类型可以调用指针方法呢?这是因为golang在编译底层做了优化，当可以拿到地址的时候，系统会自动的加上寻址符号-取指针-最终会被优化成这个样子：var b Brid = Brid{"lqm"};(&b).Fly()
	var b = Brid{"lqm"}
	b.Fly()
}
