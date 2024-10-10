package main

import (
	"log"
	"os/exec"
)

func main() {
	p := "/web/wwwroot/mlxy.burstedgold.com"
	data, err := exec.Command("sh", "/root/shellscript/svn_update2.sh", p).Output()
	if err != nil {
		log.Print("qwert11111 = ", err, string(data))
		return
	}

	log.Print("1111111 = ", string(data))
}
