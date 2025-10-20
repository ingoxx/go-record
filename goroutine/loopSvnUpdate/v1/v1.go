package main

import (
	"context"
	"encoding/json"
	"github.com/emicklei/go-restful/v3/log"
	"github.com/ingoxx/go-record/goroutine/loopSvnUpdate/v1/config"
	"io"
	"os"
	"os/exec"
	"sync"
	"time"
)

type ProjectsData struct {
	Name string
}

type TaskData struct {
	Projects []ProjectsData `json:"projects"`
}

type LoopSvnUp struct {
}

func (l *LoopSvnUp) runScript(t *ProjectsData) {
	var cancel context.CancelFunc
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*config.TimeOUT)
	defer cancel()

	cmd := exec.CommandContext(ctx, "bash", config.Script, t.Name)
	if err := cmd.Run(); err != nil {
		log.Printf("[ERROR] fail to run script, error msg '%s'\n", err.Error())
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
	var wg sync.WaitGroup

	wg.Add(config.WorkerCount)
	for i := 0; i < config.WorkerCount; i++ {
		go l.worker(ctx, tasks, &wg)
	}

	if err := l.task(tasks); err != nil {
		log.Printf("[ERROR] fail to create task, error msg '%s'\n", err.Error())
		cancel()
		return
	}

	<-stop
	cancel()
	close(tasks)
	wg.Wait()
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
			log.Printf("[WARNNING] queue was full, need to wait.")
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
	var t = LoopSvnUp{}
	t.Run()
}
