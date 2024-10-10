package main

import "fmt"

// detail01: range关键字遍历map数据是随机的，不是顺序的

// detail02: 不做任何处理直接通过range关键字遍历修改map的数据，会产生不同的情况，正式由于detail01的原因，因此不能直接通过遍历在原来的map直接插入数据

func main() {
	m := map[string]int{
		"lxb": 31,
		"lqm": 30}

	for k, v := range m {
		m[k+"2"] = v + 10
	}

	//会出现len(m)不一定是4的情况
	fmt.Println(len(m))

	//确实要修改只能是通过副本去修改,如下
	copy_m := make(map[string]int, len(m))

	for k, v := range m {
		copy_m[k] = v
	}

	for k, v := range copy_m {
		m[k+k] = v + 10
	}

}
