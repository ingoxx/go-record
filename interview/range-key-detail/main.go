package main

import "fmt"

// range的细节
func main() {
	s := []string{"hello"}
	// 不会出现死循环，range只是帮我们for循环的index去掉而已，原理是在循环之前就已经产生了副本，如下所示
	for range s {
		s = append(s, "lxb")
		fmt.Println(s)
	}

	// 原理
	// copy_s := s
	// for i := 0; i < len(copy_s); i++ {
	// 	s = append(s, "lxb")
	// 	fmt.Println(s)
	// }
}
