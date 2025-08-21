package main

import "fmt"

// slice原理
func main() {
	s := make([]int, 0)
	s2 := []int{1, 2, 3, 4, 5, 6, 7}
	s = append(s, s2...)
	fmt.Println(s, len(s), cap(s))

	s3 := []byte{22, 33, 44, 55, 66, 77, 88, 99, 11}

	var s6 []byte
	var chunks = make([][]byte, 5)
	for len(s3) >= 2 {
		s3, s6 = s3[2:], s3[:2]
		chunks = append(chunks, s6)
		fmt.Println("s6 = ", s6)
	}

	if len(s3) > 0 {
		chunks = append(chunks, s3)
	}

	s5 := len(s3)
	fmt.Println(s5, s3)
	fmt.Println(chunks)
}
