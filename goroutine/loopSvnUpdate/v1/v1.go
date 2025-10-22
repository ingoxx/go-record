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
	var cancel context.CancelFunc
	ctx, cancel := context.WithTimeout(context.Background(), config.TimeOUT)
	defer cancel()

	cmd := exec.CommandContext(ctx, "bash", config.Script, t.Name)
	if err := cmd.Run(); err != nil {
		log.Printf("[ERROR] fail to run script, project '%s' error msg '%s'\n", t.Name, err.Error())
	} else {
		log.Printf("[INFO] %s update successfully\n", t.Name)
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
	ctx, cancel := context.WithCancel(context.Background())

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	var wg sync.WaitGroup
	wg.Add(config.WorkerCount)

	for i := 0; i < config.WorkerCount; i++ {
		go l.worker(ctx, tasks, &wg)
	}

	for {
		select {
		case <-stop:
			cancel()
			close(tasks)
			wg.Wait()
			log.Println("exit ok")
			return
		case <-time.After(config.LoopSleep):
			if err := l.task(tasks); err != nil {
				log.Printf("[ERROR] fail to create task, error msg '%s'\n", err.Error())
				cancel()
				return
			}
		}
	}
}

func (l *LoopSvnUp) task(t chan *ProjectsData) error {
	td, err := l.parse()
	if err != nil {
		return err
	}

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
	log.Println("version v1.0.15")
	var t = LoopSvnUp{}
	t.Run()
}
