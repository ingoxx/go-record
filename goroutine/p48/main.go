package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

func source(c chan<- int32) {
	ra, rb := rand.Int31(), rand.Intn(3)+1
	// Sleep 1s, 2s or 3s.
	time.Sleep(time.Duration(rb) * time.Second)
	select {
	case c <- ra:
	default:
		log.Print("here")
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// The capacity should be at least 1.
	c := make(chan int32, 1)
	for i := 0; i < 5; i++ {
		go source(c)
	}
	rnd := <-c // only the first response is used
	fmt.Println(rnd)
}
