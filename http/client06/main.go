package main

import (
	"log"
	"net/http"
	"runtime"
	"sync"
	"time"
)

var (
	wg    sync.WaitGroup
	limit = make(chan struct{}, 2)
)

func main() {

	links := []string{
		"https://github.com/kubernetes/kubernetes",
		"https://github.com/kubernetes/kubernetes",
		"https://github.com/kubernetes/kubernetes",
		"https://github.com/kubernetes/kubernetes",
		"https://github.com/kubernetes/kubernetes",
		"https://github.com/kubernetes/kubernetes",
	}

	client := &http.Client{
		Transport: &http.Transport{},
		Timeout:   time.Duration(5) * time.Second,
	}

	for _, url := range links {
		limit <- struct{}{}
		wg.Add(1)
		go request(url, client)
	}

	wg.Wait()

}

func request(url string, client *http.Client) {
	log.Print("goroutine number = ", runtime.NumGoroutine())

	defer wg.Done()
	_, err := client.Get(url)
	if err != nil {
		log.Println("We could not reach:", url, err)
	} else {
		log.Println("Success reaching the website:", url)
	}
	<-limit
}
