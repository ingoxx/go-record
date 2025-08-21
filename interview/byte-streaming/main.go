package main

import (
	"fmt"
	"os"
)

func main() {
	//size := 5
	//var chunks = make([][]byte, size)

	//var chunk []byte
	s := "aabb cc ddv cc dd cc dd 234 asda"
	b := []byte(s)

	//for len(b) >= size {
	//	chunk, b = b[:size], b[size:]
	//	chunks = append(chunks, chunk)
	//}
	//
	//if len(b) > 0 {
	//	chunks = append(chunks, b[:])
	//}
	//
	chunks := chunkSlice(b, 1024)

	fmt.Println(chunks)

	file := "C:/Users/Administrator/Desktop/aa.txt"
	wi, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return
	}

	defer wi.Close()

	for _, chunk := range chunks {
		fmt.Println(string(chunk))
		_, err := wi.WriteString(string(chunk))
		if err != nil {
			return
		}
	}
}

// 将切片分割成统一的chunks块
func chunkSlice(slice []byte, chunkSize int) [][]byte {
	var chunks [][]byte

	for {
		if len(slice) == 0 {
			break
		}

		// necessary check to avoid slicing beyond
		// slice capacity
		if len(slice) < chunkSize {
			chunkSize = len(slice)
		}

		chunks = append(chunks, slice[0:chunkSize])
		slice = slice[chunkSize:]
	}

	return chunks
}
