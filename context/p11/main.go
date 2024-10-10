package main

import (
	"fmt"
	"log"
	"runtime"
	"time"
)

func f1(c chan string) {

	run()
	c <- "10 seconds passed"

}

func main() {
	c1 := make(chan string)
	go f1(c1)

	select {
	case msg1 := <-c1:
		fmt.Println(msg1)
	case <-time.After(3 * time.Second):
		fmt.Println("Timeout!")
	}

	for {
		log.Print("gn=", runtime.NumGoroutine())
		time.Sleep(time.Second)
	}
}

func run() {
	for {
		time.Sleep(time.Second * 30)

	}
}
