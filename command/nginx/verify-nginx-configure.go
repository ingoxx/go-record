package main

import (
	"fmt"
	"os"
)

func main() {
	//cmd := exec.Command("nginx", "-t")
	//err := cmd.Run()
	//if err != nil {
	//	// 如果命令执行出错，err.Error()会包含退出状态码
	//	var exitError *exec.ExitError
	//	if errors.As(err, &exitError) {
	//		// 通过ExitError的ExitCode方法获取具体的退出状态码
	//		ws := exitError.Sys().(syscall.WaitStatus)
	//		fmt.Printf("nginx -t failed with exit code: %d\n", ws.ExitStatus())
	//	}
	//	return
	//}

	err := os.Remove("C:\\Users\\Administrator\\Desktop\\sshd_config")

	fmt.Println(err == nil)
}
