package main

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"log"
	"os"
)

func main() {
	sign := "C:/Users/Administrator/Desktop/projects.json"
	user := "lxb"
	of, err := os.Open(sign)
	if err != nil {
		log.Print("projects.json not exists, esg = ", err)
		return
	}

	b, err := io.ReadAll(of)
	if err != nil {
		log.Print("read projects.json file, esg = ", err)
		return
	}

	sign = string(b)

	key := user + sign
	h := sha1.New()
	h.Write([]byte(key))
	token := hex.EncodeToString(h.Sum(nil))

	log.Print(token)
}
