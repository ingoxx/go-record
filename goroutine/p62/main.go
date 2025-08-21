package main

import (
	"log"
	"math/rand"
	"time"
)

func main() {

	rand.Seed(time.Now().UnixNano())

	c1 := make(chan int)
	c2 := make(chan int)
	// c3 := make(chan int)

	for i := 0; i < 5; i++ {
		time.Sleep(time.Duration(rand.Intn(3)+1) * time.Second)
		go func() {
			select {
			case c1 <- 1:
				<-c2
				log.Print("send c1")
			case <-c2:
			}
		}()
	}

	go func() {
		for {
			select {
			case <-c1:
				log.Print("recv c1")
				return
			default:
			}
		}
	}()

	time.Sleep(time.Second * 5)
}
