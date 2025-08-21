package main

import "fmt"

// 判断字符串字符是否都一样
func main() {
	str := "eeeee"
	b := []byte(str)
	if run(b) {
		fmt.Println("true")
	} else {
		fmt.Println("false")
	}
}

func run(b []byte) bool {
	for i := 1; i < len(b); i++ {
		index := i - 1
		t := index + 1
		if t == len(b)-1 {
			if b[t] != b[index] {
				return false
			}
		} else {
			for len(b) > t {
				if b[index] != b[t] {
					return false
				}
				t++
			}
		}

	}
	return true
}
