package main

import (
	"context"
	"encoding/json"
	"github.com/ingoxx/go-record/goroutine/loopSvnUpdate/v1/config"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"sync/atomic"
	"time"
)

type ProjectsData struct {
	Name string `json:"name"`
}

type TaskData struct {
	Projects []ProjectsData `json:"projects"`
}

type LoopSvnUp struct {
	totalData      int64
	totalProcessed int64
}

func (l *LoopSvnUp) runScript(t *ProjectsData) {
	log.Println("processed ", l.totalProcessed)
	var cancel context.CancelFunc
	ctx, cancel := context.WithTimeout(context.Background(), config.TimeOUT)
	defer cancel()

	cmd := exec.CommandContext(ctx, "bash", config.Script, t.Name)
	if err := cmd.Run(); err != nil {
		log.Printf("[ERROR] fail to run script, project '%s' error msg '%s'\n", t.Name, err.Error())
	} else {
		log.Printf("[INFO] %s update successfully\n", t.Name)
		atomic.AddInt64(&l.totalProcessed, 1)
	}
}

func (l *LoopSvnUp) worker(ctx context.Context, tasks <-chan *ProjectsData, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case t, ok := <-tasks:
			if !ok {
				return
			}
			l.runScript(t)
		}
	}
}

func (l *LoopSvnUp) Run() {
	tasks := make(chan *ProjectsData, config.QueueSize)
	stop := make(chan os.Signal, 1)
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	wg.Add(config.WorkerCount)
	for i := 0; i < config.WorkerCount; i++ {
		go l.worker(ctx, tasks, &wg)
	}

	go func() {
		for {
			select {
			case <-time.After(config.LoopSleep):
				if l.totalData == l.totalProcessed {
					atomic.AddInt64(&l.totalProcessed, 0)
					if err := l.task(tasks); err != nil {
						log.Printf("[ERROR] fail to create task, error msg '%s'\n", err.Error())
						cancel()
						return
					}
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	l.stop(cancel, tasks, stop)

	wg.Wait()

	log.Println("exist ok")
}

func (l *LoopSvnUp) stop(cancel context.CancelFunc, tasks chan *ProjectsData, s chan os.Signal) {
	signal.Notify(s, os.Interrupt)
	<-s
	cancel()
	close(tasks)
}

func (l *LoopSvnUp) task(t chan *ProjectsData) error {
	td, err := l.parse()
	if err != nil {
		return err
	}

	l.totalData = int64(len(td.Projects))
	l.totalProcessed = l.totalData

	for _, v := range td.Projects {
		select {
		case t <- &v:
		case <-time.After(config.EnqueueTimeout):
			log.Printf("[WARNNING] the queue is full and needs to wait for other goroutines to complete.")
		}
	}

	return nil
}

func (l *LoopSvnUp) readJsonFile() ([]byte, error) {
	of, err := os.Open(config.ProjectJsonFile)
	if err != nil {
		return nil, err
	}

	b, err := io.ReadAll(of)
	if err != nil {
		return nil, err
	}

	return b, nil

}

func (l *LoopSvnUp) parse() (TaskData, error) {
	var jd TaskData
	b, err := l.readJsonFile()
	if err != nil {
		return jd, err
	}

	if err := json.Unmarshal(b, &jd); err != nil {
		return jd, err
	}

	return jd, nil
}

func main() {
	log.Println("version v1.0.8")
	var t = LoopSvnUp{}
	t.Run()
}
