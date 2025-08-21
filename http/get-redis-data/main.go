package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type CpuLoad struct {
	Timestamp int64   `json:"timestamp"`
	Load1     float64 `json:"load1"`
	Load5     float64 `json:"load5"`
	Load15    float64 `json:"load15"`
	IP        string  `json:"ip"`
}

type VChartData struct {
	Columns []string         `json:"columns"`
	Rows    []map[string]any `json:"rows"`
}

func getTimePoints(days int) []int64 {
	points := []int64{}
	now := time.Now()
	location := now.Location()

	for i := days; i >= 1; i-- {
		day := now.AddDate(0, 0, -i)
		t := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, location)
		points = append(points, t.Unix())
	}
	points = append(points, now.Unix())
	return points
}

func queryCpuLoadBySegments(ctx context.Context, rdb *redis.Client, ip string, segments []int64) (map[string][]CpuLoad, error) {
	key := fmt.Sprintf("cpu_loads_%s", ip)
	values, err := rdb.LRange(key, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	segmentData := make(map[string][]CpuLoad)
	for _, val := range values {
		var load CpuLoad
		if err := json.Unmarshal([]byte(val), &load); err != nil {
			continue
		}
		for i := 0; i < len(segments)-1; i++ {
			if load.Timestamp >= segments[i] && load.Timestamp < segments[i+1] {
				label := time.Unix(segments[i], 0).Format("2006-01-02 15:04")
				segmentData[label] = append(segmentData[label], load)
				break
			}
		}
	}
	return segmentData, nil
}

func formatForVCharts(data map[string][]CpuLoad) VChartData {
	rows := []map[string]any{}
	for label, group := range data {
		if len(group) == 0 {
			continue
		}
		// 计算平均值（可选：也可选最新一条）
		sum := 0.0
		for _, item := range group {
			sum += item.Load1
		}
		avg := sum / float64(len(group))
		rows = append(rows, map[string]any{
			"时间":    label,
			"load1": fmt.Sprintf("%.2f", avg),
		})
	}

	// 排序
	sort.Slice(rows, func(i, j int) bool {
		return rows[i]["时间"].(string) < rows[j]["时间"].(string)
	})

	return VChartData{
		Columns: []string{"时间", "load1"},
		Rows:    rows,
	}
}

func main() {
	r := gin.Default()
	rdb := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6378",
		Password: "chatai",
		DB:       8,
	})

	r.GET("/api/load", func(c *gin.Context) {
		ip := c.Query("ip")
		if ip == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ip missing"})
			return
		}
		daysStr := c.DefaultQuery("days", "3")
		days, err := strconv.Atoi(daysStr)
		if err != nil || days < 1 || days > 30 {
			days = 3
		}
		segments := getTimePoints(days)
		data, err := queryCpuLoadBySegments(c, rdb, ip, segments)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, formatForVCharts(data))
	})

	r.Run(":8080")
}
