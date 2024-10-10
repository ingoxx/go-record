package main

import (
	"bufio"
	"fmt"
	"os/exec"
)

func main() {
	cmd := exec.Command("sh", "/root/shellscript/test.sh")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if err := cmd.Start(); err != nil {
		fmt.Println(err.Error())
		return
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if err := cmd.Wait(); err != nil {
		fmt.Println(err.Error())
		return
	}
}
