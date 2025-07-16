package main

import (
	"fmt"
	"github.com/ingoxx/Golang-practise/zip/zip-format"
	"log"
)

func main() {
	//src := "C:\\Users\\Administrator\\Desktop\\update"
	//dst := "C:\\Users\\Administrator\\Desktop"
	src := "C:\\Users\\Administrator\\Desktop\\update.zip"
	dst := "C:\\Users\\Administrator\\Desktop\\test"
	zp := zip_format.NewZip(src, dst)
	err := zp.UnZipFile()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("ok")
}
