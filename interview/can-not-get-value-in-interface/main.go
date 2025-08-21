package main

import "fmt"

// golang为什么无法取出接口内值的地址？
// 因为接口内存储的是值的副本，不是引用，而地址只能针对引用生效，因此golang无法取出接口内值的地址。
// 因为接口内的值是不可寻址的，只有由它所引用的值才可以取地址。实际上，接口内部存在两个值，一个是具体类型的值，另一个是实现了该接口的类型的值。因此，只有引用的具体类型值才可以取地址，而接口值本身不可寻址。
// 代码示例：

type Animal interface {
	Speak() string
}

type Dog struct {
	Name string
}

func (d Dog) Speak() string {
	return "Woof"
}

type myData struct {
	name string
}

func main() {
	data := myData{"demo"}
	var dataInterface interface{} = data
	pointer := &dataInterface //只有由它所引用的值才可以取地址-如当前
	fmt.Println(pointer)      // 取出接口的地址，打印结果：&0x82023e00

	value := *pointer
	fmt.Println(value) // 打印接口内值，打印结果：{demo}

	valuePointer := &value
	fmt.Println(*valuePointer) // 打印接口内值的地址，无法获取，panic: runtime errors: invalid memory address or nil pointer dereference

	// ----------------情形1-------------------
	m := map[string]string{
		"name": "lxb",
	}
	// _ = &m["name"] //map的value本身就是无法寻址
	fmt.Println(m)

	// ----------------情形2-------------------
	var a Animal = Dog{"eric"}
	_ = &a
	fmt.Println("a addr = ", &a)
}
