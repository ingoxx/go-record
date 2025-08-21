package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

type Config struct {
	Project []string `json:"project"`
	Limit   int      `json:"limit"`
}

func main() {
	var config Config
	file := "C:/Users/Administrator/Desktop/projects.json"
	of, err := os.Open(file)
	if err != nil {
		log.Print(err)
		return
	}

	data, err := io.ReadAll(of)
	if err != nil {
		log.Print(err)
		return
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Print(err)
		return
	}

	log.Print(config.Project)

}
