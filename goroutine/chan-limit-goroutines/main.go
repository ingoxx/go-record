package main

import (
	"log"
	"math/rand"
	"time"
)

type Customer struct{ id int }
type Bar chan Customer

func (bar Bar) ServeCustomer(c Customer) {
	log.Print("++ customer#", c.id, " starts drinking")
	time.Sleep(time.Second * time.Duration(3+rand.Intn(16)))
	log.Print("-- customer#", c.id, " leaves the bar")
	<-bar // leaves the bar and save a space
}

func main() {
	rand.Seed(time.Now().UnixNano()) // needed before Go 1.20

	// The bar can serve most 10 customers
	// at the same time.
	bar24x7 := make(Bar, 10)
	for customerId := 0; ; customerId++ {
		// time.Sleep(time.Second * 2)
		customer := Customer{customerId}
		// Wait to enter the bar.
		bar24x7 <- customer
		go bar24x7.ServeCustomer(customer)
	}

}
