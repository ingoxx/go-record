package main

import (
	"log"
	"math/rand"
	"time"
)

func main() {

	rand.Seed(time.Now().UnixNano())
	start := 0
	exit := true
	c := make(chan int, 1)

	go func() {
		data := <-c
		exit = false
		log.Print("get data = ", data)
	}()

	for exit {
		if rand.Intn(50)+1 == 21 {
			c <- 1
		}
		start++
		log.Print("try = ", start)
		time.Sleep(time.Second)
	}
}
