package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func main() {
	path := "C:\\Users\\Administrator\\Desktop\\update\\111"
	file := "C:\\Users\\Administrator\\Desktop\\update\\sql.iso"
	fn, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return
	}

	defer func(fn *os.File) {
		err := fn.Close()
		if err != nil {
			return
		}
	}(fn)

	files, err := os.ReadDir(path)
	if err != nil {
		return
	}

	var prefix = "split_"

	for k, file := range files {

		var rb = make([]byte, 10485760)

		fmt.Println(file.Name())

		fileName := fmt.Sprintf("%s%d", prefix, k)
		f, err := os.Open(filepath.Join(path, fileName))
		if err != nil {
			fmt.Println(err)
			return
		}

		for {
			rn, err := f.Read(rb)
			if err == io.EOF {
				break
			}

			if err != nil {
				fmt.Println(err)
				return
			}

			_, err = fn.Write(rb[:rn])
			if err != nil {
				fmt.Println(err)
				return
			}
		}

	}
}
