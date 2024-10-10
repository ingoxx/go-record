package main

import "fmt"

// 反转字符串
func main() {
	s3 := "我爱你, 中国"
	s4 := []int32(s3)
	for i := 1; i <= len(s4)/2; i++ {
		index := i
		offset := index - 1
		s4[offset], s4[len(s4)-index] = s4[len(s4)-index], s4[offset]
	}
	fmt.Println(string(s4))
}
