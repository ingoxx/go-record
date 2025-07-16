package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// 限制goroutine数量
var (
	wg    sync.WaitGroup
	c     = make(chan int, 20)
	root  = "D:/project/gin/src/github.com/ingoxx/Golang-practise/.git"
	total = 0
	count = make(chan int)
	stop  = make(chan struct{})
)

func main() {

	go func() {
		for {
			select {
			case <-stop:
				return
			case <-count:
				total++
				fmt.Println(total)
			}
		}
	}()

	for i := 0; i < 10000; i++ {
		c <- 1
		wg.Add(1)
		go recv(c)
	}

	wg.Wait()
	stop <- struct{}{}
	fmt.Println("finished!!!", total)
}

func recv(c chan int) {

	defer wg.Done()

	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			count <- 1
			<-c
		}
		return nil
	})
	// time.Sleep(time.Second * 3)

}
