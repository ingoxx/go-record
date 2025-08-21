package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	item := 50

	stream := func() <-chan int {
		st := make(chan int)

		go func() {
			defer close(st)
			for i := 0; i < item; i++ {
				//time.Sleep(time.Duration(rand.Intn(3)+1) * time.Second)
				st <- i
			}
		}()

		return st
	}

	run(stream())
}

func run(in <-chan int) {
	for {
		select {
		case v, ok := <-in:
			if !ok {
				return
			}
			fmt.Println(v)
		default:
		}
	}
}
