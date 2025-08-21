package main

import (
	"fmt"
	"time"
)

func main() {
	var block chan struct{}
	var c = make(chan bool)
	go func() {
		time.Sleep(time.Duration(1) * time.Second)
		fmt.Println(<-c)
	}()
	c <- true
	<-block
}
