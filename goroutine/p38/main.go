package main

import (
	"log"
	"math/rand"
	"runtime"
	"time"
)

var (
	workers = make(chan int, 10)
)

type Seat int
type Bar chan Seat

// 最多同时10个客户，但还有大量客户需要goroutine
func (bar Bar) ServeCustomerAtSeat(c int, seat Seat) {
	log.Print("gn111 = ", runtime.NumGoroutine())
	log.Print("++ customer#", c, " drinks at seat#", seat)
	time.Sleep(time.Second * time.Duration(2+rand.Intn(6)))
	log.Print("-- customer#", c, " frees seat#", seat)
	bar <- seat // free seat and leave the bar
	// <-workers
}

func main() {

	rand.Seed(time.Now().UnixNano())

	bar24x7 := make(Bar, 10)
	for seatId := 0; seatId < cap(bar24x7); seatId++ {
		bar24x7 <- Seat(seatId)
	}

	for customerId := 0; ; customerId++ {
		log.Print("gn2222 = ", runtime.NumGoroutine())
		// time.Sleep(time.Second)
		// Need a seat to serve next customer.
		seat := <-bar24x7
		// workers <- 10
		go bar24x7.ServeCustomerAtSeat(customerId, seat)
	}

}
