package main

import (
	"log"
	"math/rand"
	"time"
)

type Seat int
type Bar chan Seat

func (bar Bar) ServeCustomerAtSeat(c int, seat Seat) {
	log.Print("++ customer#", c, " drinks at seat#", seat)
	time.Sleep(time.Second * time.Duration(2+rand.Intn(6)))
	log.Print("-- customer#", c, " frees seat#", seat)
	bar <- seat // free seat and leave the bar
}

func main() {
	rand.Seed(time.Now().UnixNano()) // needed before Go 1.20

	bar24x7 := make(Bar, 10)
	for seatId := 0; seatId < cap(bar24x7); seatId++ {
		bar24x7 <- Seat(seatId)
	}
	go func() {
		for {
			log.Print(">>>", len(bar24x7))
			if len(bar24x7) == 0 {
				log.Print("null")
				return
			}
		}
	}()

	for customerId := 0; customerId < 30; customerId++ {
		time.Sleep(time.Second)
		// Need a seat to serve next customer.
		seat := <-bar24x7
		go bar24x7.ServeCustomerAtSeat(customerId, seat)
	}

	for {
		time.Sleep(time.Second)
	}
}
