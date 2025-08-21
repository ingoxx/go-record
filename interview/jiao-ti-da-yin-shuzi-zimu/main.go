package main

import (
	"fmt"
)

func main() {
	s1 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	s2 := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}
	var done = make(chan int)
	var recs1 = make(chan int)
	var recs2 = make(chan string)
	var f int

	go func() {
		for i := 0; i < len(s1); i++ {
			recs1 <- s1[i]
		}
		done <- 1
	}()

	go func() {
		for i := 0; i < len(s2); i++ {
			recs2 <- s2[i]
		}
		done <- 1
	}()

	for {
		select {
		case num := <-recs1:
			fmt.Println(num)
		case letter := <-recs2:
			fmt.Println(letter)
		case <-done:
			if f == 1 {
				return
			}
			f++
		}
	}

}
