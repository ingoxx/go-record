package main

import (
	"fmt"
	"os"
)

func main() {
	b, _ := os.ReadFile("C:/Users/Administrator/Desktop/create_ag_9639.sql")
	fmt.Println(string(b))
}
