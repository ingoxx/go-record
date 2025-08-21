package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Result struct {
	Response interface{}
	Err      error
}

type Task interface {
	Execute() Result
}

type WorkerPool struct {
	wg     *sync.WaitGroup
	mu     *sync.Mutex
	cond   *sync.Cond
	once   *sync.Once
	quit   chan interface{}
	taskCh chan Task
	ctx    context.Context

	counter int
}

func NewWorkerPool(taskLimit int, ctx context.Context) *WorkerPool {
	return &WorkerPool{
		mu:     new(sync.Mutex),
		quit:   make(chan interface{}),
		taskCh: make(chan Task, taskLimit),
		cond:   sync.NewCond(new(sync.Mutex)),
		once:   new(sync.Once),
		wg:     new(sync.WaitGroup),
		ctx:    ctx,
	}
}

func (w *WorkerPool) Start(n int) *sync.WaitGroup {
	w.wg.Add(n)
	for i := 0; i < n; i++ {
		go w.loop()
	}
	fmt.Printf("started %d workers\n", n)
	return w.wg
}

func (w *WorkerPool) AddTask(tasks ...Task) {
	for _, task := range tasks {
		select {
		case <-w.quit:
			return
		case w.taskCh <- task:
			w.cond.L.Lock()
			w.cond.Broadcast()
			fmt.Println("received task", task)
			w.cond.L.Unlock()
		}
	}
}

func (w *WorkerPool) loop() {
	defer w.wg.Done()
	for {
		select {
		case <-w.ctx.Done():
			return
		case task, ok := <-w.taskCh:
			if !ok {
				return
			}
			res := task.Execute()
			w.mu.Lock()
			w.counter++
			w.mu.Unlock()
			fmt.Println("task:", res)
		default:
			w.cond.L.Lock()
			fmt.Println("taskCh full, wait consume.", len(w.taskCh))
			w.cond.Wait()
			w.cond.L.Unlock()
		}
	}
}

func (w *WorkerPool) Stop() {
	w.once.Do(func() {
		close(w.taskCh)
		w.cond.L.Lock()
		w.cond.Broadcast()
		w.cond.L.Unlock()
	})
}

type DelayTask struct{}

func (d *DelayTask) Execute() Result {
	time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
	return Result{
		Response: "done",
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wp := NewWorkerPool(100, ctx)

	numWorkers := 5
	job := wp.Start(numWorkers)
	go func() {
		for i := 0; i < 100; i++ {
			wp.AddTask(&DelayTask{})
		}
		wp.Stop()
	}()

	job.Wait()
	fmt.Println("exiting", wp.counter)
}
