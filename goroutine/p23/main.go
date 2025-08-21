package main

import (
	"fmt"
	"time"
)

func main() {
	c := make(chan int) // an unbuffered channel
	go func(ch chan<- int, x int) {
		fmt.Println("阻塞发送")
		time.Sleep(time.Second * 5)
		// <-ch    // fails to compile
		// Send the value and block until the result is received.
		ch <- x * x // 9 is sent
	}(c, 3)

	done := make(chan struct{})
	go func(ch <-chan int) {
		// Block until 9 is received.
		fmt.Println("阻塞接收")
		time.Sleep(time.Second * 7)
		n := <-ch
		fmt.Println(n) // 9
		// ch <- 123   // fails to compile
		time.Sleep(time.Second * 3)
		done <- struct{}{}
	}(c)
	// Block here until a value is received by
	// the channel "done".
	<-done

	fmt.Println("bye")
}
