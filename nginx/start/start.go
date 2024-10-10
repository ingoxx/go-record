package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func main() {

	cmd := exec.Command("/usr/sbin/nginx", "-c", "/etc/nginx/nginx.conf")
	var out strings.Builder
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		log.Fatalf("fail to start nginx: %v", err)
	}

	fmt.Println("NGINX START...")

	for {

	}
}
