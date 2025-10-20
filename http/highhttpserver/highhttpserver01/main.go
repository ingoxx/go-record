package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"time"
)

// Task 表示入队的任务（这里简单包含 ID 和 Payload）
type Task struct {
	ID      int64  `json:"id"`
	Payload string `json:"payload"`
	Created time.Time
}

// Configuration / tuning
const (
	QueueSize         = 1000                   // 队列缓冲大小（可根据内存 / QPS 调整）
	WorkerCount       = 50                     // worker 数量（并发处理数）
	EnqueueTimeout    = 100 * time.Millisecond // 当队列已满时，尝试放入队列的超时时间（实现 backpressure）
	ServerAddr        = ":8080"
	SimulatorClients  = 200 // 模拟器并发客户端数
	SimulatorRequests = 500 // 每个模拟客户端发送的请求数
)

// Metrics
var (
	totalEnqueued  int64
	totalProcessed int64
	totalRejected  int64
)

// worker 处理函数（可以替换为具体业务，如发起下游 HTTP 请求、数据库操作等）
func worker(ctx context.Context, id int, tasks <-chan *Task, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Printf("[worker-%02d] started\n", id)
	for {
		select {
		case <-ctx.Done():
			log.Printf("[worker-%02d] shutting down\n", id)
			return
		case t, ok := <-tasks:
			if !ok {
				log.Printf("[worker-%02d] tasks channel closed\n", id)
				return
			}
			// 模拟处理耗时
			processTask(t, id)
			atomic.AddInt64(&totalProcessed, 1)
		}
	}
}

func processTask(t *Task, workerID int) {
	// 这里的处理逻辑可以替换为实际工作（如发起 HTTP 请求等）
	// 为了演示，我们随机延迟 10-200ms
	d := time.Duration(10+rand.Intn(190)) * time.Millisecond
	time.Sleep(d)
	if rand.Intn(1000) == 0 { // 少量概率打印详细日志
		log.Printf("[worker-%02d] processed task id=%d payload=%q (took=%s)\n", workerID, t.ID, t.Payload, d)
	}
}

func enqueueHandler(tasks chan<- *Task) http.HandlerFunc {
	var idCounter int64
	return func(w http.ResponseWriter, r *http.Request) {
		// 读取 body（可根据需要解析 JSON）
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		id := atomic.AddInt64(&idCounter, 1)
		task := &Task{
			ID:      id,
			Payload: string(body),
			Created: time.Now(),
		}

		// 尝试放入队列（有超时），超时则返回 503 表示被拒绝（backpressure）
		select {
		case tasks <- task:
			atomic.AddInt64(&totalEnqueued, 1)
			w.WriteHeader(http.StatusAccepted)
			fmt.Fprintf(w, "enqueued id=%d\n", task.ID)
		case <-time.After(EnqueueTimeout):
			atomic.AddInt64(&totalRejected, 1)
			http.Error(w, "queue full, try later", http.StatusServiceUnavailable)
		}
	}
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	s := map[string]int64{
		"enqueued":  atomic.LoadInt64(&totalEnqueued),
		"processed": atomic.LoadInt64(&totalProcessed),
		"rejected":  atomic.LoadInt64(&totalRejected),
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(s)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// 有界队列 channel
	tasks := make(chan *Task, QueueSize)

	// 启动 worker pool
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(WorkerCount)
	for i := 0; i < WorkerCount; i++ {
		go worker(ctx, i+1, tasks, &wg)
	}

	// HTTP Server
	mux := http.NewServeMux()
	mux.HandleFunc("/enqueue", enqueueHandler(tasks))
	mux.HandleFunc("/stats", statsHandler)

	srv := &http.Server{
		Addr:    ServerAddr,
		Handler: mux,
	}

	// Graceful shutdown on SIGINT/SIGTERM
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go func() {
		log.Printf("server listening on %s\n", ServerAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen error: %v", err)
		}
	}()

	// 启动模拟器（可注释掉，如果你想用外部工具压测）
	go simulateLoad()

	<-stop
	log.Println("shutting down server...")

	// 1) 停止接收新的连接
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}

	// 2) 通知 workers 关闭并等待其完成已取到的任务
	cancel()
	// 关闭 tasks channel 让 worker 遍历结束
	close(tasks)
	wg.Wait()

	log.Println("all workers stopped")
	log.Printf("final stats: enqueued=%d processed=%d rejected=%d\n",
		atomic.LoadInt64(&totalEnqueued),
		atomic.LoadInt64(&totalProcessed),
		atomic.LoadInt64(&totalRejected),
	)
}

// simulateLoad 会并发地向 /enqueue 发送大量请求，用于压力测试演示
func simulateLoad() {
	time.Sleep(500 * time.Millisecond) // 等服务起来
	log.Println("simulator: start sending requests...")

	var clientsWG sync.WaitGroup
	clientsWG.Add(SimulatorClients)
	client := &http.Client{Timeout: 3 * time.Second}

	for c := 0; c < SimulatorClients; c++ {
		go func(clientID int) {
			defer clientsWG.Done()
			for i := 0; i < SimulatorRequests; i++ {
				// 模拟不同负载与请求体大小
				payload := fmt.Sprintf("client=%d seq=%d rnd=%d", clientID, i, rand.Intn(1e6))
				req, _ := http.NewRequest(http.MethodPost, "http://127.0.0.1"+ServerAddr+"/enqueue", io.NopCloser(
					strReader(payload),
				))
				resp, err := client.Do(req)
				if err != nil {
					// 超时或连接错误
					// log.Printf("sim client %d req %d error: %v", clientID, i, err)
					time.Sleep(5 * time.Millisecond)
					continue
				}
				_ = resp.Body.Close()
				// 少量节流，避免瞬间把客户端本身打爆
				time.Sleep(time.Duration(rand.Intn(5)) * time.Millisecond)
			}
		}(c)
	}

	clientsWG.Wait()
	log.Println("simulator: finished sending requests")
}

// strReader 是 io.Reader 辅助
func strReader(s string) io.Reader {
	return &stringReader{s, 0}
}

type stringReader struct {
	s string
	i int64
}

func (r *stringReader) Read(p []byte) (n int, err error) {
	if r.i >= int64(len(r.s)) {
		return 0, io.EOF
	}
	n = copy(p, r.s[r.i:])
	r.i += int64(n)
	return
}
