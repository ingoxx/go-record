package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func main() {

	// filename := "/Users/liaoxuanbiao/Downloads/log/20221216/2022121618.log.gz"
	filename := "C:/Users/Administrator/Desktop/log/20221216/2022121618.log.gz"
	wfile := filepath.Join(filepath.Dir(filename), "2022121618.log")
	of, err := os.Create(wfile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	gzipfile, err := os.Open(filename)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	reader, err := gzip.NewReader(gzipfile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer reader.Close()

	if _, err = io.Copy(of, reader); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
