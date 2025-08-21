package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

var (
	send    = make(chan string)
	wg      sync.WaitGroup
	finish  = make(chan string)
	worker  = 8
	date    = "20240718"
	suffix  = ".bak"
	project = "dbbak"
	disk    = "F"
)

func main() {
	now := time.Now()
	cmd := "rsync"
	host := fmt.Sprintf("10.0.135.127::nas/%s/db-1/", project)
	dirs := []string{
		"F:\\dbbak\\EC2AMAZ-HGIUIOA",
	}
	dirsLen := len(dirs)
	var over = make(chan struct{}, dirsLen)
	ctx := context.Background()

	wg.Add(worker)
	for i := 0; i < worker; i++ {
		go func(ctx context.Context) {
			defer wg.Done()
			ctx1, cancel := context.WithTimeout(ctx, time.Second*time.Duration(7200))
			defer cancel()
			for file := range send {
				if err := exec.CommandContext(ctx1, cmd, "-av", "--ignore-errors", file, host).Run(); err != nil {
					log.Println("exec cmd err >>> ", err, file)
				} else {
					fmt.Printf("successful %s \n", file)
				}
			}
		}(ctx)
	}

	go func() {
		for v := range finish {
			loopDir(v)
			over <- struct{}{}
		}
	}()

	go func() {
		for {
			select {
			case <-over:
				if dirsLen == 1 {
					close(send)
					return
				}
				dirsLen--
			}
		}
	}()

	for _, dir := range dirs {
		finish <- dir
	}

	close(finish)
	wg.Wait()

	close(over)

	fmt.Printf("传输总共耗时 >>> %s", time.Now().Sub(now))
}

func loopDir(dir string) {
	readDir, err := os.ReadDir(dir)
	if err == nil {
		for _, v := range readDir {
			if v.IsDir() {
				loopDir(filepath.Join(dir, v.Name()))
			} else {
				if strings.HasSuffix(filepath.Join(dir, v.Name()), suffix) && isMatch(filepath.Join(dir, v.Name())) {
					rf := replace(filepath.Join(dir, v.Name()))
					send <- rf
				}
			}
		}
	}
}

func replace(path string) string {
	file := strings.Split(path, ":")
	pf := filepath.Join(fmt.Sprintf("/cygdrive/%s/", disk), file[1])
	return strings.ReplaceAll(pf, "\\", "/")
}

func isMatch(file string) bool {
	p := fmt.Sprintf("%s.*.bak$", date)
	re := regexp.MustCompile(p)
	re.FindString("db")
	return re.MatchString(file)
}
