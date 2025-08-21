package main

import (
	"fmt"
	"strings"
)

func main() {
	str := "aa dd dd	ab4"
	d := strings.Split(str, ",")
	fmt.Println(d, len(d), d[len(d)-1])
}
