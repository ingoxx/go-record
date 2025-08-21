package main

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

// M个接收者，一个发送者”情况的变种：关闭请求由第三方goroutine发出

func main() {
	rand.Seed(time.Now().UnixNano()) // before Go 1.20
	log.SetFlags(0)

	// ...
	const Max = 100000
	const NumReceivers = 100
	const NumThirdParties = 5

	wgReceivers := sync.WaitGroup{}
	wgReceivers.Add(NumReceivers)

	// ...
	dataCh := make(chan int)
	closing := make(chan struct{}) // signal channel
	closed := make(chan struct{})

	// The stop function can be called
	// multiple times safely.
	stop := func() {
		select {
		case closing <- struct{}{}:
			<-closed
			// log.Println("ttttttttttttttttttttttttttttttttttttttttt")
		case <-closed:
		}
	}

	// some third-party goroutines
	for i := 0; i < NumThirdParties; i++ {
		go func() {
			r := 1 + rand.Intn(3)
			time.Sleep(time.Duration(r) * time.Second)
			stop()
		}()
	}

	// the sender
	go func() {
		defer func() {
			close(closed)
			close(dataCh)
		}()

		for {
			select {
			case <-closing:
				log.Println("bbbcccccccccccccccccccccccccccccccccccccccc")
				return
			default:
			}

			select {
			case <-closing:
				log.Println("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
				return
			case dataCh <- rand.Intn(Max):
			}
		}
	}()

	// receivers
	for i := 0; i < NumReceivers; i++ {
		go func() {
			defer wgReceivers.Done()

			for value := range dataCh {
				log.Println(value)
			}
		}()
	}

	wgReceivers.Wait()
}
