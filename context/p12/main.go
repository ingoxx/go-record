package main

import (
	"context"
	"log"
	"os/exec"
	"time"
)

func main() {

	timeout := 3

	ctx := context.Background()

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "sleep", "1")
	if err := cmd.Run(); err != nil {
		log.Print("timeout")
	} else {
		log.Print("finished")
	}

}
