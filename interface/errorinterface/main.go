package main

type Pointer struct{}

func main() {
	p := newPointer()
	if p == nil {
		println("p is nil")
	} else {
		println("p is not nil")
	}

}

// detail01: golang是根据一个值的前八个字节是否为0值来判断是否是nil/空值
// detail02: interface是一个16个字节长的结构体，首8字节是类型字段，后8个字节是数据指针
// detail03: nil指针 != nil接口，nil指针占用的内存大小取决于指针数据类型

func newPointer() interface{} {
	//这里给c为Pointer结构体指针类型赋了nil值，必然会导致接口的前8个字节不等于空值
	var c *Pointer = nil
	return c
}

// 这里是涉及到，interface经常会被滥用，如需要返回nil，应该明确返回具体值，如上例子，只需把返回值interface改成*Pointer，或者直接return nil
