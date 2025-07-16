package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

func countFiles(root string, fileCount chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	count := 0
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error accessing %q: %v\n", path, err)
			return err
		}
		if !info.Mode().IsRegular() {
			// Ignore non-regular files like directories or symlinks
			return nil
		}
		if info.Name()[0] == '.' {
			// Ignore hidden files and directories
			return nil
		}
		count++
		return nil
	})
	fileCount <- count
}

func main() {
	dir := "/Users/liaoxuanbiao/project/golang/src/github.com/ingoxx/Golang-practise/.git"
	numThreads := 4
	fileCount := make(chan int, numThreads)
	var wg sync.WaitGroup
	wg.Add(numThreads)
	for i := 0; i < numThreads; i++ {
		go countFiles(dir, fileCount, &wg)
	}
	go func() {
		wg.Wait()
		close(fileCount)
	}()
	total := 0
	for count := range fileCount {
		total += count
	}
	fmt.Printf("Number of files in %s: %d\n", dir, total)
}
