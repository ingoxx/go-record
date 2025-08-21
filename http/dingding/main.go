package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type WarnInfo struct {
	webHook string
	data    string
}

func NewWarnInfo(data string) *WarnInfo {
	return &WarnInfo{
		webHook: "https://oapi.dingtalk.com/robot/send?access_token=4797ec430d2b74acbfaa084960ca389d884068a0a1d6115ad10c4ac7ffabc395",
		data:    data,
	}
}

func (w *WarnInfo) SendWarningInfo() {
	data := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]interface{}{
			"content": w.data,
		},
		"at": map[string]interface{}{
			"atMobiles": "15889709122",
			"isAtAll":   false,
		},
	}

	w.send(data)
}

func (w *WarnInfo) send(data map[string]interface{}) error {
	b, _ := json.Marshal(data)
	body := bytes.NewBuffer(b)
	resp, _ := http.Post(w.webHook, "application/json", body)
	io.ReadAll(resp.Body)

	return nil
}

func main() {
	var data = fmt.Sprintf("告警信息\n终端操作： docker info")
	NewWarnInfo(data).SendWarningInfo()
}
