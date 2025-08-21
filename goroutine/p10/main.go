package main

import "fmt"

func main() {
	// g := Gen(10, 23)
	// c := Sq(g)
	// fmt.Println(<-c)
	// fmt.Println(<-c)

	// for v := range c {
	// 	fmt.Println(v)
	// }

	for v := range Sq(Gen(10, 23)) {
		fmt.Println(v)
	}
}

func Gen(num ...int) <-chan int {
	c := make(chan int)
	go func() {
		for _, v := range num {
			c <- v
		}
		close(c)
	}()

	return c
}

func Sq(c <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for v := range c {
			out <- v * v
		}
		close(out)
	}()

	return out

}
