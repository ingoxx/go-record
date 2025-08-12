package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

func main() {
	mvFile()
}

func rename() {
	path := "C:\\Users\\Administrator\\Desktop\\img3"
	dir, err := os.ReadDir(path)
	if err != nil {
		log.Fatalln(err)
	}
	var count = 1000
	for _, v := range dir {
		if !v.IsDir() {
			of := filepath.Join(path, v.Name())
			stat, err := os.Stat(of)
			if err != nil {
				log.Printf("[ERROR] fail to get size %s\n", v.Name())
				continue
			}
			if stat.Size() > 2048 {
				nf := filepath.Join(path, fmt.Sprintf("%s.png", strconv.Itoa(count)))
				if err := os.Rename(of, nf); err != nil {
					log.Printf("[ERROR] fail to rename %s\n", v.Name())
				}
				count += 1
			}
		}
	}

	log.Println("finished")
}

func mvFile() {
	tarPath := "C:\\Users\\Administrator\\Desktop\\profile3"
	path := "C:\\Users\\Administrator\\Desktop\\img3"
	dir, err := os.ReadDir(path)
	if err != nil {
		log.Fatalln(err)
	}
	var count = 1001
	for _, v := range dir {
		if !v.IsDir() {
			of := filepath.Join(path, v.Name())
			stat, err := os.Stat(of)
			if err != nil {
				log.Printf("[ERROR] fail to get size %s\n", v.Name())
				continue
			}
			if stat.Size() > 2048 {
				nf := filepath.Join(tarPath, fmt.Sprintf("%s.png", strconv.Itoa(count)))
				create, _ := os.Create(nf)
				read, _ := os.OpenFile(of, os.O_RDONLY, 0777)
				_, err := io.Copy(create, read)
				if err != nil {
					log.Printf("[ERROR] fail to copy %s\n", v.Name())
				}
				count += 1
			}
		}
	}

	log.Println("finished")
}
