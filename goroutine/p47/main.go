package main

import (
	"log"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

// 使用带缓冲的channel来限制goroutine数量，（会重复创建goroutine，但问题不大，Golang的协程创建消耗不大）
type Customer struct{ id int }
type Bar chan Customer

var wg sync.WaitGroup

func (bar Bar) ServeCustomer(c Customer) {
	defer wg.Done()
	log.Print("++ customer#", c.id, " starts drinking")
	log.Print("gn = ", runtime.NumGoroutine())
	time.Sleep(time.Second * time.Duration(1+rand.Intn(3)))
	// time.Sleep(time.Second * 30)
	log.Print("-- customer#", c.id, " leaves the bar ", len(bar))
	<-bar // leaves the bar and save a space
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// The bar can serve most 10 customers
	// at the same time.
	bar24x7 := make(Bar, 10)
	for customerId := 0; customerId <= 100; customerId++ {
		//time.Sleep(time.Second / 7)
		customer := Customer{customerId}
		select {
		case bar24x7 <- customer:
			wg.Add(1) // try to enter the bar
			go bar24x7.ServeCustomer(customer)
		default:
			log.Print("customer#", customerId, " goes elsewhere")
			// <-bar24x7
		}
	}

	wg.Wait()
}
