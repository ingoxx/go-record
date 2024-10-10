package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {

	// create a channel for work "tasks"
	ch := make(chan int)

	wg := sync.WaitGroup{}

	// start the workers
	for t := 0; t < 100; t++ {
		wg.Add(1)
		go saveToDB(ch, &wg)
	}

	// push the lines to the queue channel for processing
	for i := 0; i < 100000; i++ {
		ch <- i
		fmt.Println("1111111111111111111111")
	}

	// this will cause the workers to stop and exit their receive loop
	close(ch)

	// make sure they all exit
	wg.Wait()
}

func saveToDB(ch chan int, wg *sync.WaitGroup) {
	fmt.Println("gn1111111 = ", runtime.NumGoroutine())
	// cnosume a line
	for line := range ch {
		// do work
		fmt.Println(line)
		fmt.Println("gn222222 = ", runtime.NumGoroutine())
	}
	// we've exited the loop when the dispatcher closed the channel,
	// so now we can just signal the workGroup we're done
	wg.Done()
}
