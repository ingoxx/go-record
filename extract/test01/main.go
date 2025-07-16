package main

import (
	"log"

	"github.com/ingoxx/Golang-practise/extract"
)

func main() {
	filename := "C:/Users/Administrator/Desktop/log/20221216/2022121618.log.gz"
	unGz := extract.NewUngz(filename)
	err := unGz.UngzFile()
	if err != nil {
		log.Print(err)
		return
	}

	log.Print("finished")
}
