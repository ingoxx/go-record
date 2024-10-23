package main

import (
	"bytes"
	"fmt"
	"github.com/mitchellh/go-ps"
	"k8s.io/klog/v2"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

var (
	pid = "/var/run/nginx.pid"
)

func isRunning() bool {
	log.Println("isRunning")
	processes, err := ps.Processes()
	if err != nil {
		klog.ErrorS(err, "unexpected errors obtaining process list")
	}
	for _, p := range processes {
		fmt.Println(p.Executable())
		if p.Executable() == "nginx" {
			return true
		}
	}

	return false
}

func checkNginxConfig() (bool, string, error) {
	cmd := exec.Command("nginx", "-t")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	// 执行命令
	err := cmd.Run()

	// 检查命令执行结果
	if err != nil {
		// 如果执行失败，err包含了退出状态码
		if exitError, ok := err.(*exec.ExitError); ok {
			return false, strings.TrimSpace(stderr.String()), nil
		} else {
			// 其他类型的错误
			return false, "", exitError
		}
	}

	// 如果没有错误，说明配置文件是正确的
	return true, "Configuration file test is successful", nil
}

func reload() error {
	output, err := exec.Command("cat", pid).Output()
	if err != nil {
		return err
	}

	ngxPid, err := strconv.Atoi(strings.TrimSpace(string(output)))
	if err != nil {
		return err
	}

	if err = syscall.Kill(ngxPid, syscall.SIGHUP); err != nil {
		klog.ErrorS(err, "failed to reload nginx, rollback in progress")
		return err
	}

	log.Println("reload nginx successfully")

	return nil
}

func overwrite(src, dst string) error {
	readFile, err := os.ReadFile(src)
	if err == nil {
		if err := os.WriteFile(dst, readFile, 0644); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	_, s, err := checkNginxConfig()
	if err != nil {
		log.Fatalln("check errors: ", err)
	}

	log.Println(s)

	if err := reload(); err != nil {
		log.Fatalln("reload errors: ", err)
	}

	running := isRunning()
	log.Println("isRunning >>> ", running)
}
