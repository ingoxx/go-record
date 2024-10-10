package main

import (
	"fmt"
)

func main() {
	c := make(chan int, 3)
	fmt.Println(len(c))
	for i := 0; i < 3; i++ {
		c <- i
	}

	fmt.Println(len(c))

}
