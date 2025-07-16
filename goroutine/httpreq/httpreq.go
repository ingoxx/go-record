package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

func main() {
	urls := []string{
		"https://www.google.com/",
		"https://docs.aws.amazon.com/zh_cn/inspector/latest/user/scanning-ec2.html#deep-inspection",
		"https://docs.aws.amazon.com/aws-managed-policy/latest/reference/AmazonSSMManagedInstanceCore.html",
		"https://github.com/ingoxx/chat",
		"https://stackoverflow.com/",
		"https://docs.aws.amazon.com/systems-manager/latest/userguide/agent-install-rhel.html",
		"https://www.baidu.com/",
		"https://www.qq.com/",
	}

	ctx, cancel := context.WithCancel(context.Background())

	var workers = 10
	var countFinishedWorker int

	req := newRequest(ctx, workers)
	req.wg.Add(1)

	req.start()

	go func() {
		for {
			select {
			case <-req.ctx.Done():
				return
			case <-req.finished:
				countFinishedWorker++
				if countFinishedWorker == workers {
					close(req.result)
					close(req.err)
					cancel()
					return
				}
			}
		}
	}()

	go func() {
		defer req.wg.Done()
		for {
			select {
			case <-req.ctx.Done():
				return
			case v, ok := <-req.result:
				if !ok {
					return
				}
				fmt.Println(v)
			case e, ok := <-req.err:
				if !ok {
					return
				}
				fmt.Println("err >>> ", e)
			}
		}
	}()

	go func() {
		for _, v := range urls {
			req.task <- v
		}
		req.stop()
	}()

	req.wg.Wait()

}

type result struct {
	url  string
	resp string
}

//type task func() *result

type request struct {
	ctx      context.Context
	task     chan string
	result   chan *result
	workers  int
	finished chan struct{}
	wg       *sync.WaitGroup
	lock     *sync.Mutex
	err      chan error
}

func newRequest(ctx context.Context, workers int) *request {
	return &request{
		ctx:      ctx,
		task:     make(chan string),
		result:   make(chan *result),
		finished: make(chan struct{}),
		wg:       new(sync.WaitGroup),
		lock:     new(sync.Mutex),
		workers:  workers,
		err:      make(chan error),
	}
}

func (r *request) start() {
	for i := 0; i < r.workers; i++ {
		go r.work()
	}
}

func (r *request) run(u string) *result {
	r.lock.Lock()
	defer r.lock.Unlock()

	ctx1, cancel1 := context.WithTimeout(context.Background(), time.Second*time.Duration(2))
	defer cancel1()

	var req = new(result)
	var done = make(chan *http.Response)

	go func() {
		fb, err := os.ReadFile("")
		if err != nil {
			r.err <- err
			return
		}
		var params = map[string]interface{}{
			"file": fb,
		}
		jb, err := json.Marshal(&params)
		if err != nil {
			r.err <- err
			return
		}

		body := bytes.NewReader(jb)

		req, err := http.NewRequest(http.MethodPost, u, body)
		if err != nil {
			r.err <- err
			return
		}
		req = req.WithContext(ctx1)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			r.err <- err
			return
		}
		defer resp.Body.Close()
		done <- resp
	}()

	select {
	case <-ctx1.Done():
		req.url = u
		req.resp = "time out"
	case v := <-done:
		req.url = v.Request.URL.String()
		req.resp = v.Request.Host
	}

	return req
}

func (r *request) stop() {
	close(r.task)
}

func (r *request) work() {
	for {
		select {
		case <-r.ctx.Done():
			return
		case v, ok := <-r.task:
			if !ok {
				r.finished <- struct{}{}
				return
			}
			r.result <- r.run(v)
		}
	}
}
