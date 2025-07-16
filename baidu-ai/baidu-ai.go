package main

import (
	"encoding/json"
	"fmt"
	"github.com/ingoxx/Golang-practise/http/newHttp"
	"io"
	"log"
	"net/http"
	"time"
)

var (
	client = &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:          100,
			MaxIdleConnsPerHost:   100,
			MaxConnsPerHost:       100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   35 * time.Second,
			ExpectContinueTimeout: 35 * time.Second,
		},
		Timeout: time.Duration(35) * time.Second,
	}
)

type Msg struct {
	Role    string
	Content string
}

type Resp struct {
	Result string `json:"result"`
}

func main() {
	chat()
}

func getToken() string {
	key := "b6xsGr4XQo0NgKmgBWUpXGTP"
	secret := "5ksZiuOEU2rbdkZNm3gKdgUsvH9zX3U6"
	url := fmt.Sprintf("https://aip.baidubce.com/oauth/2.0/token?client_id=%s&client_secret=%s&grant_type=client_credentials", key, secret)
	var headers = make(map[string]interface{})
	var data = make(map[string]interface{})
	var resp = make(map[string]interface{})
	var b = make([]byte, 0)
	headers["Content-Type"] = "application/json"
	headers["Accept"] = "application/json"
	req := newHttp.HttpReq(url, "aa", data, headers, false, 5)
	body, err := req.GET(client)
	if err != nil {
		log.Fatalln(err)
	}

	defer body.Close()

	b, err = io.ReadAll(body)
	if err != nil {
		log.Fatalln(err)
	}

	err = json.Unmarshal(b, &resp)
	if err != nil {
		log.Fatalln(err)
	}

	return resp["access_token"].(string)
}

func chat() {
	var headers = make(map[string]interface{})
	var data = make(map[string][]*Msg)
	var resp Resp
	var msg = make([]*Msg, 0)
	url := "https://aip.baidubce.com/rpc/2.0/ai_custom/v1/wenxinworkshop/chat/eb-instant?access_token=" + getToken()
	headers["Content-Type"] = "application/json"
	headers["Transfer-Encoding"] = "chunked"
	d1 := &Msg{
		Role:    "user",
		Content: "介绍下你自己",
	}

	msg = append(msg, d1)
	data["messages"] = msg
	b, _ := json.Marshal(data)
	req := newHttp.HttpReq(url, string(b), nil, headers, true, 30)
	body, err := req.POST(client)
	if err != nil {
		log.Fatalln(err)
	}

	defer body.Close()

	b2, err := io.ReadAll(body)
	if err != nil {
		log.Fatalln(err)
	}

	json.Unmarshal(b2, &resp)
	fmt.Println(string(b2))

}
