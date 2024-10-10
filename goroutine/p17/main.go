package main

import (
	"fmt"
	"time"
)

func main() {
	messages := make(chan int)

	// Use this channel to follow the execution status
	// of our goroutines :D
	done := make(chan bool)

	// go func() {
	// 	time.Sleep(time.Second * 3)
	// 	messages <- 1
	// 	done <- true
	// }()
	// go func() {
	// 	time.Sleep(time.Second * 2)
	// 	messages <- 2
	// 	done <- true
	// }()
	// go func() {
	// 	time.Sleep(time.Second * 1)
	// 	messages <- 3
	// 	done <- true
	// }()

	for i := 1; i < 4; i++ {
		go func(i int) {
			time.Sleep(time.Second)
			messages <- i
			done <- true
		}(i)
	}

	go func() {
		for i := range messages {
			fmt.Println(i)
		}
	}()

	for i := 0; i < 3; i++ {
		<-done
		fmt.Println(11111)
	}
}
