package main

import "fmt"

type Person struct {
	Name string
}

func (p Person) GetName() {
	fmt.Println("name =", p.Name)
}

type Man struct {
	Person
}

type Kind interface {
	GetName()
}

func GetName(p Kind) {
	p.GetName()
}

func main() {
	m := Man{
		Person{Name: "lxb"},
	}

	GetName(m)

	// ------------------------------------
	//golang是根据变量的前8个字节是否为零值来判断是否为nil
	var x interface{} = nil //空的接口赋值给空的接口，16个字节都为零,所以等于nil
	var y *int = nil        //空的指针对接口进行赋值会导致前8个字节不为空，所以不等于nil，接口是前8个自己指向类型，后8个字节指向数据
	interfaceNil(x)
	interfaceNil(y)
}

func interfaceNil(x interface{}) {
	if x == nil {
		fmt.Println("empty interface")
		return
	}
	fmt.Println("non-empty interface")
}
