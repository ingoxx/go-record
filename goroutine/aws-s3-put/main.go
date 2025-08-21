package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var (
	//发送的任务的协程数
	limitChan = make(chan struct{}, 20)
	WorkChan  = make(chan string)
	wgSend    sync.WaitGroup
	wgRec     sync.WaitGroup
	//单位秒
	timeOut    = 600
	putSrcPath = flag.String("src", "", "need to upload path")
)

func LoopDir(root string, limit chan struct{}, finished bool) {
	fd, err := os.ReadDir(root)
	if err == nil {
		for _, file := range fd {
			if strings.Contains(filepath.Join(root, file.Name()), "sbl_db") {
				if file.Name() == "MGLog" || file.Name() == "LOG" {
					continue
				}

				if file.IsDir() {
					select {
					case limit <- struct{}{}:
						wgSend.Add(1)
						go LoopDir(filepath.Join(root, file.Name()), limit, false)
					default:
						LoopDir(filepath.Join(root, file.Name()), limit, true)
					}
				} else {
					WorkChan <- filepath.Join(root, file.Name())
				}
			}
		}
	}

	if !finished {
		wgSend.Done()
		<-limit
	}
}

func RunCmd(file string, ctx context.Context) (err error) {
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, time.Second*time.Duration(timeOut))

	defer cancel()

	cmd := fmt.Sprintf("s3://db-backup-huawen/truco/%s", filepath.Base(file))

	c := exec.CommandContext(ctx, "aws", "s3", "cp", file, cmd)
	if err = c.Run(); err != nil {
		return
	}

	return
}

func main() {
	flag.Parse()

	if flag.NFlag() != 1 {
		log.Fatalln(flag.ErrHelp.Error())
	}

	start := time.Now()
	rootDir := *putSrcPath

	//接收的任务的协程数
	const recWork int = 20
	ctx := context.Background()

	wgRec.Add(recWork)
	for range [recWork]struct{}{} {
		go func() {
			defer wgRec.Done()
			for file := range WorkChan {
				if err := RunCmd(file, ctx); err == nil {
					log.Printf("%s succeed to upload aws s3", filepath.Base(file))
				} else {
					log.Printf("%s failed to upload aws s3, esg >>> %s", filepath.Base(file), err.Error())
				}
			}
		}()
	}

	LoopDir(rootDir, limitChan, true)

	wgSend.Wait() //这里是为了等待遍历完所有目录，然后关闭WorkChan

	close(WorkChan)

	wgRec.Wait() //这里是为了等待所有文件都上传完

	fmt.Println("done, cost time >>>", time.Since(start))
}
