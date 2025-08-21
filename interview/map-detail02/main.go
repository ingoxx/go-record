package main

import "fmt"

//map的一些细节
// detail：如果是嵌套比较深的map，建议使用指针
type MyStruct struct {
	Name string
}

func main() {
	m := make(map[string]*MyStruct)

	m["s01"] = &MyStruct{Name: "lxb"}

	sm, ok := m["s01"]
	if ok {
		sm.Name = "lqm" //m := make(map[string]MyStruct)如果map存的是结构体而不是结构体指针，单纯的这样修是不生效的，还得把修改好的sm.Name 重新赋值给sm，如下
		// m["s01"] = sm
	}

	//无法获取val的内存地址，它通过哈希对比随机放到不同的bucket
	// _ = &m["s01"]

	fmt.Println(m["s01"])
}
