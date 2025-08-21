package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	total = 0

	root = "C:\\Windows"
)

func main() {
	start := time.Now()

	Loop(root)

	fmt.Printf("total = %d, time = %v\n", total, time.Since(start))

}

func Loop(root string) {
	fd, err := os.ReadDir(root)
	if err == nil {
		for _, file := range fd {
			if !file.IsDir() {

				if strings.HasSuffix(file.Name(), ".sql") {
					fmt.Println(root, file.Name())
				}

				total++
			} else {
				Loop(filepath.Join(root, file.Name()))
			}
		}
	}
}
