package main

import (
	"fmt"
	"time"
)

func isNextDay(t1, t2 time.Time) bool {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()

	if y1 != y2 || m1 != m2 || d1 != d2 {
		return true
	}
	return false
}

func main() {
	t1 := time.Now()
	fmt.Println(t1.Format(time.DateTime))
	t2 := t1.Add(3 * time.Hour) // 添加24小时，模拟第二天的时间

	fmt.Println(isNextDay(t1, t2)) // 输出“true”
}
