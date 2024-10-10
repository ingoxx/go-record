package main

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
)

type Seat int
type Bar chan Seat

var (
	// wg      sync.WaitGroup
	work      = make(chan string)
	totalChan = make(chan int)
	workers   = make(chan int, 20)
	total     = 0
)

func getwork() {
	// for path := range work {
	// 	workers <- 1
	// 	go sendwork(path)
	// }

	for {
		select {
		case path := <-work:
			go sendwork(path)
		}
	}
}

func sendwork(path string) {
	fl, err := os.ReadDir(path)
	if err == nil {
		for _, file := range fl {
			if file.IsDir() {
				work <- filepath.Join(path, file.Name())
			} else {
				totalChan <- 1
			}
		}
	}
	<-workers
}

func main() {

	path := "C:/Windows/"
	// path := "C:/Users/Administrator/Desktop/test/"
	go func() {
		for {
			select {
			case <-totalChan:
				total++
			default:
				log.Print("total =", runtime.NumGoroutine())
			}
		}
	}()

	go sendwork(path)

	getwork()

	select {}

}
