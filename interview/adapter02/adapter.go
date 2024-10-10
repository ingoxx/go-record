package main

import "fmt"

//在Go语言中，适配器设计模式主要用于将一个接口转换为另一个接口，使得原本不兼容的类可以协同工作。以下是一个使用适配器模式的简单示例：
//假设我们有一个已存在的OldCalculator接口和实现，但新的业务逻辑需要使用新定义的Calculator接口

// 已有的旧计算器接口和其实现

type OldCalculator interface {
	Add(int, int) int
	Subtract(int, int) int
}

type oldCalculator struct{}

func (oc oldCalculator) Add(a, b int) int {
	return a + b
}

func (oc oldCalculator) Subtract(a, b int) int {
	return a - b
}

// 新的计算器接口，新增了一个Multiply方法

type Calculator interface {
	Add(int, int) int
	Subtract(int, int)
	Multiply(int, int) int
}

// 创建一个适配器，实现新的Calculator接口并调用旧的OldCalculator方法

type Adapter struct {
	oldCalc OldCalculator
}

func NewAdapter(oc OldCalculator) *Adapter {
	return &Adapter{oldCalc: oc}
}

func (a *Adapter) Add(x, y int) int {
	return a.oldCalc.Add(x, y)
}

func (a *Adapter) Subtract(x, y int) int {
	return a.oldCalc.Subtract(x, y)
}

// 通过适配器提供旧接口没有的Multiply方法

func (a *Adapter) Multiply(x, y int) int {
	// 假设这里实现了乘法操作
	return x * y
}

func main() {
	oldCalc := oldCalculator{}
	adapter := NewAdapter(oldCalc)

	// 现在我们可以使用新的Calculator接口与旧的OldCalculator一起工作
	sum := adapter.Add(3, 5)
	diff := adapter.Subtract(7, 2)
	product := adapter.Multiply(4, 6)

	fmt.Println("Sum:", sum)
	fmt.Println("Difference:", diff)
	fmt.Println("Product:", product)
}

//在这个例子中，Adapter扮演了适配器的角色，它实现了新的Calculator接口，并在内部使用了OldCalculator的方法。
//同时，对于OldCalculator未提供的功能（如Multiply方法），我们在适配器中进行了补充实现。
