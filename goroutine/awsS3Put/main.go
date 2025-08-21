package main

import (
	"flag"
	"fmt"
	"gopkg.in/gcfg.v1"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/ingoxx/Golang-practise/aws/s3"
)

var (
	limitChan = make(chan struct{}, 20)
	WorkChan  = make(chan string)
	wgSend    sync.WaitGroup
	wgRec     sync.WaitGroup
	iniFile   = flag.String("ini", "", "ini file path")
)

type S3Config struct {
	Key        string
	Secret     string
	Region     string
	PutSrcPath string
	Bucket     string
}

type Config struct {
	HuaWen S3Config
}

func main() {
	flag.Parse()

	if flag.NFlag() != 1 {
		log.Fatalln(flag.ErrHelp.Error())
	}

	var cfg Config
	err := gcfg.ReadFileInto(&cfg, *iniFile)
	if err != nil {
		log.Println(err)
		return
	}

	start := time.Now()
	const recWork int = 10
	root := cfg.HuaWen.PutSrcPath
	config := []string{cfg.HuaWen.Key, cfg.HuaWen.Secret, cfg.HuaWen.Region}

	so := s3.NewS3Client(config...)
	sc, err := so.S3Client()
	if err != nil {
		return
	}

	s3api := &s3.Object{
		Bucket:   cfg.HuaWen.Bucket,
		S3Client: sc,
	}

	wgRec.Add(recWork)
	for range [recWork]struct{}{} {
		go func() {
			defer wgRec.Done()
			for file := range WorkChan {

				err := s3api.PutLargeObject(file, "truco/"+filepath.Base(file))
				if err == nil {
					log.Printf("%s succeed to upload aws s3", filepath.Base(file))
				} else {
					log.Printf("%s failed to upload aws s3, esg >>> %s", filepath.Base(file), err.Error())
				}
			}
		}()
	}

	LoopDir(root, limitChan, true)

	wgSend.Wait() //这里是为了等待遍历完所有目录，然后关闭WorkChan

	close(WorkChan)

	wgRec.Wait() //这里是为了等待所有文件都上传完

	fmt.Printf("time = %v\n", time.Since(start))
}

func LoopDir(root string, limit chan struct{}, finished bool) {
	fd, err := os.ReadDir(root)
	if err == nil {
		for _, file := range fd {
			if strings.Contains(filepath.Join(root, file.Name()), "split_") {
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
