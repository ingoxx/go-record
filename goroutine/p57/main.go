package main

import (
	"log"
	"math/rand"
	"runtime"
	"strconv"
	"sync"
	"time"
)

// M 个接收者，N 个发送者，其中任何一个通过通知主持人关闭一个额外的信号通道来表示“让我们结束游戏”

func main() {
	rand.Seed(time.Now().UnixNano()) // before Go 1.20
	log.SetFlags(0)

	// ...
	const Max = 100000
	const NumReceivers = 10
	const NumSenders = 30

	wgReceivers := sync.WaitGroup{}
	wgReceivers.Add(NumReceivers)

	// ...
	dataCh := make(chan int)
	stopCh := make(chan struct{})
	// stopCh is an additional signal channel.
	// Its sender is the moderator goroutine shown
	// below, and its receivers are all senders
	// and receivers of dataCh.
	toStop := make(chan string, 1)
	//因为是多个发送者，中间人必须是有缓冲chan且cap必须为1，这样才能在第一次收到消息后第一时间close(stopCh)
	// The channel toStop is used to notify the
	// moderator to close the additional signal
	// channel (stopCh). Its senders are any senders
	// and receivers of dataCh, and its receiver is
	// the moderator goroutine shown below.
	// It must be a buffered channel.

	var stoppedBy string

	// moderator
	go func() {
		stoppedBy = <-toStop
		close(stopCh)
	}()

	// senders
	for i := 0; i < NumSenders; i++ {
		go func(id string) {
			defer func() { log.Print("send stop1111111111111111") }()
			for {
				log.Print("gn11111111 = ", runtime.NumGoroutine())
				time.Sleep(time.Duration(rand.Intn(10) * int(time.Second)))
				value := rand.Intn(Max)
				if value >= 97000 {
					// Here, the try-send operation is
					// to notify the moderator to close
					// the additional signal channel.
					select {
					case toStop <- "sender#" + id: // 30个发送goroutine，但只有最先进入到这里的goroutine才能close(stopCh)，其他的goroutine进到这里都会阻塞直到被gc回收, 所以只会打印一次send stop1111111111111111，这个就是设置管道容量为1的好处
					default:
					}
					return

				}

				// The try-receive operation here is to
				// try to exit the sender goroutine as
				// early as possible. Try-receive and
				// try-send select blocks are specially
				// optimized by the standard Go
				// compiler, so they are very efficient.
				select {
				case <-stopCh:
					return
				default:
				}

				// Even if stopCh is closed, the first
				// branch in this select block might be
				// still not selected for some loops
				// (and for ever in theory) if the send
				// to dataCh is also non-blocking. If
				// this is unacceptable, then the above
				// try-receive operation is essential.
				// 当两个case都是非阻塞的时候select case是随机的，所以才会有上面的一个try receives
				// 即使已经关闭了stopCh, 还是有可能先走第二个case
				select {
				case <-stopCh:
					return
				case dataCh <- value:
				}
			}
		}(strconv.Itoa(i))
	}

	// receivers
	for i := 0; i < NumReceivers; i++ {
		go func(id string) {
			defer func() { wgReceivers.Done(); log.Print("recv stop1111111111111111111") }()

			for {
				// Same as the sender goroutine, the
				// try-receive operation here is to
				// try to exit the receiver goroutine
				// as early as possible.
				select {
				case <-stopCh:
					return
				default:
				}

				// Even if stopCh is closed, the first
				// branch in this select block might be
				// still not selected for some loops
				// (and forever in theory) if the receive
				// from dataCh is also non-blocking. If
				// this is not acceptable, then the above
				// try-receive operation is essential.
				select {
				case <-stopCh:
					return
				case value := <-dataCh:
					if value == Max-1 {
						// Here, the same trick is
						// used to notify the moderator
						// to close the additional
						// signal channel.
						select {
						case toStop <- "receiver#" + id:
						default:
						}
						return
					}
					//do something
					log.Println(value)
				}
			}
		}(strconv.Itoa(i))
	}

	// ...
	wgReceivers.Wait()
	log.Println("stopped by", stoppedBy)
}
