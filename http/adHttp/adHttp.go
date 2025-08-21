package adhttp

import (
	"log"
	"net/http"
	"sync"
	"time"
)

type HighHttpClient struct {
	hr      *http.Request
	Timeout int
	once    sync.Once
	client  *http.Client
}

func (hhc *HighHttpClient) GET(url string) {
	hhc.once.Do(func() {
		hhc.NewRequest()
	})

	_, err := hhc.client.Get(url)
	if err != nil {
		log.Println("We could not reach:", url, err)
	} else {
		log.Println("Success reaching the website:", url)
	}
}

func (hhc *HighHttpClient) NewRequest() {
	tr := &http.Transport{}
	hhc.client = &http.Client{
		Transport: tr,
		Timeout:   time.Duration(hhc.Timeout) * time.Second,
	}
}

func NewHighHttpClient(timeout int) *HighHttpClient {
	return &HighHttpClient{
		Timeout: timeout,
	}
}
