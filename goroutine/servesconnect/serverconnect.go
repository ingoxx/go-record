package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// WebSocket 连接升级器
var upgrader = websocket.Upgrader{
	// 解决跨域问题
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 定义返回给前端的结果结构体
type ConnResult struct {
	IP      string `json:"ip"`
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"` // omitempty 表示如果字段为空，则在JSON中省略
}

// WebSocket 处理函数
func wsHandler(c *gin.Context) {
	// 1. 升级 HTTP 连接为 WebSocket 连接
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Failed to upgrade connection:", err)
		return
	}
	defer conn.Close() // 确保函数退出时关闭连接

	log.Println("Client connected:", conn.RemoteAddr())

	// 2. 进入消息监听循环
	for {
		// 读取来自客户端的消息
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			// 如果是连接关闭的错误，则正常退出
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("Client disconnected unexpectedly:", err)
			} else {
				log.Println("Client disconnected:", conn.RemoteAddr())
			}
			break
		}

		// 我们只处理文本消息
		if messageType == websocket.TextMessage {
			var ips []string
			// 解析收到的JSON数组
			if err := json.Unmarshal(p, &ips); err != nil {
				log.Println("Failed to unmarshal JSON:", err)
				// 可以给前端返回一个错误信息
				conn.WriteJSON(gin.H{"error": "Invalid JSON format. Expected an array of strings."})
				continue
			}

			log.Printf("Received IPs to check: %v\n", ips)

			// 3. 并发处理IP连接
			results := processIPs(ips)

			// 4. 将结果返回给前端
			if err := conn.WriteJSON(results); err != nil {
				log.Println("Failed to write JSON response:", err)
				break
			}
			log.Printf("Sent %d results back to client.\n", len(results))
		}
	}
}

// 并发处理IP列表的核心函数
func processIPs(ips []string) []ConnResult {
	var wg sync.WaitGroup
	// 创建一个带缓冲的channel，大小与IP数量相同，避免Goroutine阻塞
	resultsChan := make(chan ConnResult, len(ips))

	// 定义连接超时时间
	const connectionTimeout = 2 * time.Second

	for _, ip := range ips {
		wg.Add(1) // WaitGroup 计数器加1

		// 为每个IP启动一个Goroutine
		go func(ipAddr string) {
			defer wg.Done() // Goroutine结束时，计数器减1

			address := net.JoinHostPort(ipAddr, "9092")
			result := ConnResult{IP: ipAddr}

			// 尝试连接，设置超时
			conn, err := net.DialTimeout("tcp", address, connectionTimeout)
			if err != nil {
				result.Success = false
				result.Error = err.Error()
			} else {
				result.Success = true
				conn.Close() // 检查成功后立即关闭连接
			}

			// 将结果发送到channel
			resultsChan <- result
		}(ip)
	}

	// 启动一个Goroutine来等待所有任务完成，然后关闭channel
	// 这是一个非常重要的模式，可以防止主流程死锁
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// 从channel中收集所有结果
	var finalResults []ConnResult
	for res := range resultsChan {
		finalResults = append(finalResults, res)
	}

	return finalResults
}

func main() {
	router := gin.Default()

	// 设置WebSocket路由
	router.GET("/ws", wsHandler)

	// 启动服务器
	log.Println("Server started on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to run server:", err)
	}
}
