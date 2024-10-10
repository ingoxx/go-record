package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	file, err := os.Stat("aaa")
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(file)
}
