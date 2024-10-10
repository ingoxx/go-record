package main

import (
	"fmt"
	"io/fs"
	"path/filepath"
)

func main() {

	dir := "C:\\Users\\Administrator\\Desktop\\update"

	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			//fileName := filepath.Join(path, info.Name())
			fmt.Printf("full path >>>, %s, filename >>> %s\n", path, info.Name())
		}

		return nil
	})

	if err != nil {
		fmt.Println(err)
	}
}
