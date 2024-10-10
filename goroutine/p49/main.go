package main

import (
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"time"
)

func source() <-chan int32 {
	// c must be a buffered channel.
	c := make(chan int32, 1)
	go func() {
		ra, rb := rand.Int31(), rand.Intn(3)+1
		time.Sleep(time.Duration(rb) * time.Second)
		c <- ra
	}()
	return c
}

func main() {
	rand.Seed(time.Now().UnixNano())

	var rnd int32
	// Blocking here until one source responses. 这里会发生内存泄漏，有两个goroutine永远挂起
	select {
	case rnd = <-source():
	case rnd = <-source():
	case rnd = <-source():
	}
	fmt.Println(rnd)

	for {
		time.Sleep(time.Second)
		log.Print("gn = ", runtime.NumGoroutine())
	}
}
