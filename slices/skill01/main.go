package main

import (
	"fmt"
)

// 传参切片小技巧
func main() {
	//skill 01
	Name([]string{"lxb", "lqm"}...)
	//skill 02
	Name("lxb", "lqm")

	//善用命名返回值, 可自动的帮我们初始化
	s1 := []string{"lxb", "lqm"}
	s2 := Copys(s1)
	fmt.Println(s2)

}

// 不定参数
func Name(p ...string) {
	for k, v := range p {
		println(k, v)
	}
}

func Copys(src []string) (dst []string) {
	// dst = append(dst, src...)
	//或者用内置函数
	copy(dst, src)
	return dst
}
