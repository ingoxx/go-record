package main

import (
	"fmt"
	"github.com/ingoxx/Golang-practise/gzip/gzip-new"
	"log"
)

func main() {
	//src := "C:\\Users\\Administrator\\Desktop\\update"
	//dst := "C:\\Users\\Administrator\\Desktop"
	src := "C:\\Users\\Administrator\\Desktop\\update.tar.gz"
	dst := "C:\\Users\\Administrator\\Desktop\\test"
	zp := gzip_new.NewGzip(src, dst)
	err := zp.UnGzipFile()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("ok")
}
