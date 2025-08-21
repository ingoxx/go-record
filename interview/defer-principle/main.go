package main

import "log"

// Go defer关键字的实现原理：先进后出

func main() {
	res := test1(4)
	log.Print("res >>>", res)
}

func test1(i int) int {
	defer test3()

	if i > 2 {
		defer test2()
	}
	log.Print("test1")
	return i
}

//实际执行的伪代码
// func test1(i int) int {
// 	log.Print("test1")
// 	test2()
// 	test3()
// 	return i
// }

func test2() int {
	log.Print("test2")
	return 2
}

func test3() int {
	log.Print("test3")
	return 3
}
