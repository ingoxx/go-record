package main

import (
	"fmt"
	"sync"
)

type Job func()

type WorkerPool struct {
	jobs      chan Job
	waitGroup sync.WaitGroup
}

func NewWorkerPool(size int, bufferSize int) *WorkerPool {
	pool := &WorkerPool{
		jobs: make(chan Job, bufferSize),
	}

	pool.waitGroup.Add(size)

	for i := 0; i < size; i++ {
		go func() {
			defer pool.waitGroup.Done()
			for job := range pool.jobs {
				job()
			}
		}()
	}

	return pool
}

func (pool *WorkerPool) AddJob(job Job) {
	pool.jobs <- job
}

func (pool *WorkerPool) Close() {
	close(pool.jobs)
	pool.waitGroup.Wait()
}

func main() {
	numWorkers := 5
	bufferSize := 100
	pool := NewWorkerPool(numWorkers, bufferSize)

	for i := 0; i < 10; i++ {
		i := i
		pool.AddJob(func() {
			fmt.Printf("Job %d executed by worker\n", i)
		})
	}

	pool.Close()
}
