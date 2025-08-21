package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

var (
	wg        sync.WaitGroup
	jobs      = make(chan string)
	totalChan = make(chan int)
	total     = 0
)

func loopFilesWorker() error {
	for path := range jobs {
		// fmt.Println("gn0000000000000000000 = ", runtime.NumGoroutine())
		files, err := os.ReadDir(path)
		if err != nil {
			wg.Done()
			return err
		}

		for _, file := range files {
			if !file.IsDir() {
				// fmt.Println(file.Name())
				fmt.Println("gn1111111111111111111 = ", runtime.NumGoroutine())
			}
		}
		wg.Done()
	}
	return nil
}

func LoopDirsFiles(path string) error {
	// fmt.Println("gn22222222222222222 = ", runtime.NumGoroutine())
	files, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	//Add this path as a job to the workers
	//You must call it in a go routine, since if every worker is busy, then you have to wait for the channel to be free.
	wg.Add(1)
	go func() {
		jobs <- path
	}()
	for _, file := range files {
		if file.IsDir() {
			//Recursively go further in the tree
			LoopDirsFiles(filepath.Join(path, file.Name()))
		}
	}
	return nil
}

func main() {
	// path := "C:/Windows/"
	path := "C:/Users/Administrator/Desktop/test/"
	//Start as many workers you want, now 10 workers
	for w := 1; w <= 10; w++ {
		go loopFilesWorker()
	}
	//Start the recursion
	LoopDirsFiles(path)
	wg.Wait()
}
