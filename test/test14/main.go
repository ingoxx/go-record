package main

import "fmt"

func test1(n int) {
	if n > 2 {
		n--
		test1(n)
		fmt.Println(111111)
	}
	fmt.Println("test1 n = ", n)
}

func test4(n int) int {
	if n == 1 {
		return 3
	} else {
		return 2*test4(n-1) + 1
	}
}

func main() {
	test1(4)
	t4 := test4(2)
	fmt.Println(t4)
}
