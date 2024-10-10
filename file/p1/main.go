package main

import (
	"fmt"
	"os"
)

func main() {
	path := "C:/Users/Administrator/Desktop/test/"
	fl, err := os.ReadDir(path)
	if err == nil {
		for _, file := range fl {
			if file.IsDir() {
				fmt.Println("dir = ", file.Name())
			} else {
				fmt.Println("file = ", file.Name())
			}
		}
	}
}
