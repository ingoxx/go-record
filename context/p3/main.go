package main

import "fmt"

func main() {
	i := 1
LOOP:
	for {
		if i == 10 {
			break LOOP
		}
		fmt.Println(i)
		i++
	}
	i++
	goto LOOP
}
