package main

import "fmt"

// 不能在循环中修改map，因为range取map的顺序是不固定的
func main() {
	m := map[string]int{
		"name": 1,
		"addr": 2,
	}

	for k, v := range m {
		m[k+"_new"] = v
	}

	fmt.Println("m errors = ", len(m), m)

	// 解决方案，取出keys，遍历加入到map

	m1 := map[string]int{
		"name": 1,
		"addr": 2,
	}

	var sk []string

	for k := range m1 {
		sk = append(sk, k+"_new")
	}

	fmt.Println("sk = ", sk)

	for k, v := range sk {
		m1[v] = k
	}

	fmt.Println("m1 succeed = ", len(m1), m1)

}
