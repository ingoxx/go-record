package main

import "log"

func main() {
	var m1 []int
	m := make([]int, 3)
	m1 = append(m1, 2)
	m = append(m, 1)
	m = append(m, 1)
	m = append(m, 1)
	m = append(m, 1)
	log.Print(m, m1)
}
