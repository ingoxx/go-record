package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var (
	// wg        sync.WaitGroup
	childpath = make(chan string)
	sem       = make(chan bool, 50)
	totalChan = make(chan bool)
	total     = 0
)

func Task(path string) {
	select {
	case childpath <- path:
	case sem <- true:
		go RecvPath(path)
	}
}

func SendPath(path string) {
	fl, err := os.ReadDir(path)
	if err == nil {
		for _, file := range fl {
			if file.IsDir() {
				Task(filepath.Join(path, file.Name()))
				SendPath(filepath.Join(path, file.Name()))
			}
		}
	}
}

func RecvPath(path string) {
	defer func() { <-sem }()

	for {
		fl, err := os.ReadDir(path)
		if err == nil {
			for _, file := range fl {
				if !file.IsDir() {
					totalChan <- true
					fmt.Println("gn = ", runtime.NumGoroutine())
				}
			}
		} else {
			fmt.Println("err = ", err)
		}

		<-childpath

	}
}

// func Count() {
// 	for {
// 		select {
// 		case <-totalChan:
// 			total++
// 			fmt.Println("total = ", total)
// 		}
// 	}
// }

func main() {
	// path := "C:/Windows/"
	path := "C:/Users/Administrator/Desktop/test/"
	start := time.Now()
	// go Count()
	SendPath(path)
	time.Sleep(time.Second * 50)
	fmt.Printf("total = %d, cost time = %v", total, time.Since(start))
}
