package main

// 快速实现接口的方法
type myInterface interface {
	myMethod1()
	myMethod2()
}

type myStruct struct{}

// myMethod1 implements myInterface
func (*myStruct) myMethod1() {
	panic("unimplemented")
}

// myMethod2 implements myInterface
func (*myStruct) myMethod2() {
	panic("unimplemented")
}

func main() {
	var _ myInterface = (*myStruct)(nil)

}
