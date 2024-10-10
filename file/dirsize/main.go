package main

import (
	"fmt"
	"os"
	"runtime"
	"time"
)

var (
	childpath = make(chan string)
	workers   = make(chan bool, 20)
	finished  = make(chan bool)
	totalchan = make(chan bool)
	done      = 1
	total     = 0
)

func Work(path string, size int) {
	// fmt.Printf("协程数: %d\n", runtime.NumGoroutine())
	defer func() {
		finished <- true
	}()

	fl, err := os.ReadDir(path)
	if err == nil {
		for _, file := range fl {
			if file.IsDir() {
				childpath <- path + file.Name() + "/"
			} else {
				totalchan <- true
				fmt.Printf("协程数: %d\n", runtime.NumGoroutine())
			}
		}
	}

	<-workers
}

func Run(path string, size int) {
	go Work(path, size)

	for {
		select {
		case t := <-childpath:
			done++
			workers <- true
			go Work(t, size)
		case <-finished:
			done--
			if done == 0 {
				return
			}
		case <-totalchan:
			total++
		}
	}
}

func main() {
	// path := "C:/Users/Administrator/Desktop/test/"
	path := "C:/Windows/"
	size := 0

	// flag.StringVar(&path, "path", "", "目录名")
	// flag.IntVar(&size, "size", 0, "要查找的文件大小")

	// flag.Parse()
	// if flag.NFlag() != 2 {
	// 	fmt.Println(flag.ErrHelp.Error() + ", input -h for help")
	// 	return
	// }

	start := time.Now()
	Run(path, size)
	fmt.Printf("total= %d, 耗时: %v\n", total, time.Since(start))
}
