package main

import (
	"log"
	"math/rand"
	"time"
)

func main() {

	rand.Seed(time.Now().UnixNano())

	for {
		log.Print("try = ", rand.Intn(50)+1)
		time.Sleep(time.Second)
	}
}
