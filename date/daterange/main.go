package main

import (
	"fmt"
	"time"
)

func main() {
	days := 3
	var points []int64
	now := time.Now()
	location := now.Location()
	for i := days - 1; i >= 0; i-- {
		day := now.AddDate(0, 0, -i)
		dayStart := time.Date(day.Year(), day.Month(), day.Day(), 23, 59, 59, 59, location)
		points = append(points, dayStart.Unix())
	}

	fmt.Println(points)
}
