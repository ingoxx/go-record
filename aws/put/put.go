package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Lxb921006/Golang-practise/aws/s3"
)

var (
	wg         sync.WaitGroup
	put        = make(chan string)
	limit      = make(chan struct{}, 20)
	stop       = make(chan struct{})
	recvn      = make(chan struct{})
	n1         = 0
	n2         = 0
	iniFile    = flag.String("ini", "", "ini file path")
	section    = flag.String("section", "", "ini section")
	region     = flag.String("region", "", "aws region")
	putSrcPath = flag.String("src", "", "upload file")
)

func main() {
	flag.Parse()

	if flag.NFlag() != 4 {
		log.Fatalln(flag.ErrHelp.Error())
	}

	wg.Add(20)

	config := []string{*iniFile, *section, *region}
	s3api := &s3.Object{
		Bucket: "db-backup-huawen",
		S3Sess: s3.NewS3Sess(config...),
	}

	for range [20]struct{}{} {
		go func() {
			defer wg.Done()
			for {
				select {
				case <-stop:
					return
				default:
				}

				select {
				case <-stop:
					return
				case v := <-put:
					err := s3api.PutObject(v, "truco/"+filepath.Base(v))
					if err == nil {
						log.Printf("%s succeed to upload aws s3", filepath.Base(v))
					} else {
						log.Printf("%s failed to upload aws s3, esg >>> %s", filepath.Base(v), err.Error())
					}
				}
			}
		}()
	}

	fd, err := os.ReadDir(*putSrcPath)
	if err != nil {
		log.Fatalln(err.Error())
	}

	for _, file := range fd {
		if strings.HasPrefix(file.Name(), "sbl_") {
			limit <- struct{}{}
			go TargetDir(filepath.Join(*putSrcPath, file.Name()), true, limit)
			n1++
		}
	}

	go func() {
		for {
			select {
			case <-recvn:
				n2++
				if n2 == n1 {
					close(stop)
					return
				}
			default:
			}
		}
	}()

	wg.Wait()
}

func TargetDir(root string, exit bool, l <-chan struct{}) {
	fd, err := os.ReadDir(root)
	if err == nil {
		for _, v := range fd {
			if v.Name() == "MGLog" || v.Name() == "LOG" {
				continue
			}
			if v.IsDir() {
				TargetDir(filepath.Join(root, v.Name()), false, l)
			} else {
				put <- filepath.Join(root, v.Name())
			}
		}
	} else {
		log.Printf("dir %s not exists", root)
	}

	if exit {
		recvn <- struct{}{}
		<-l
	}
}
