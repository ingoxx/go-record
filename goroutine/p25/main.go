package main

import "fmt"

func main() {
	cc := make(chan int)

	recv := func(i int) {
		for {
			fmt.Println(<-cc, "退不出来")
			i++
			cc <- i
		}
	}

	go recv(1)
	go recv(2)

	cc <- 1

	var c chan bool // nil
	<-c             // blocking here for ever

	fmt.Println("end")
}
