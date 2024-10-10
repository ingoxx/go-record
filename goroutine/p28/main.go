package main

import "fmt"

func main() {
	var c chan struct{} // nil

	// fmt.Println(c)

	select {
	case c <- struct{}{}: // blocking operation
	case <-c: // blocking operation
		fmt.Println("recv")
	default:
		fmt.Println("Go here.")
	}
}
