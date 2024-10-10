package main

import "fmt"

func main() {
	c := make(chan int, 2)
	go func() {
		for i := 0; i < 10; i++ {
			c <- i
		}
		// close(c)
	}()

	// for {
	// 	v, k := <-c
	// 	if !k {
	// 		break
	// 	}
	// 	fmt.Println(v)

	// }

	for v := range c {
		fmt.Println(v)
	}

	fmt.Println("main finished")
}
