package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	//当没有使用goroutine时，无缓冲通道必须要有接受才能发送，否则发成死锁

	// i := 1
	c := make(chan int)
	select {
	case <-c:
		log.Print("recv")
	default:
		log.Print("default")
	}

	// c <- 1
	// c <- 2
	// c <- 3
	// <-c
	// <-c
	// close(c)
	// v, k := <-c
	// fmt.Println(v, k)
	// limit := make(chan int, 5)
	// c <- 1
	// c <- 1
	// c <- 1
	// for i := 0; i < 100; i++ {
	// 	go Recv(c, limit)
	// }

	// for {
	// 	i++
	// 	c <- i
	// }
	// <-c
	// c <- 2
	// c <- 10 //这里会收阻塞，无缓冲通道在发的时候都会阻塞直到有人接收
	// fmt.Println("send finished")
}

func requestWaitTimeout() {
	c1 := make(chan int, 1)
	c2 := make(chan int)

	go Recv(c1, c2)

	select {
	case <-c1:
		log.Print("c1 缓冲")
	case <-c2:
		log.Print("c2 无缓冲")
	}
}

func Recv(c1, c2 chan int) {
	// fmt.Println("resc succeed = ", <-c)
	time.Sleep(time.Second)
	c2 <- 1
	// c <- 1
}

func Send(c chan int) {
	c <- 1
	fmt.Println("send succeed")
}
