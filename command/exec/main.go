package main

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
)

func main() {
	makeCmd := fmt.Sprintf("sh cmd.sh \"more /var/log/aaa.log\"")
	cmd := exec.Command("sh", "-c", makeCmd)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalln(err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatalln("Error getting stderr:", err)
	}

	if err = cmd.Start(); err != nil {
		log.Fatalln(err)
	}

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	if err = cmd.Wait(); err != nil {
		log.Fatalln(err)
	}
}
