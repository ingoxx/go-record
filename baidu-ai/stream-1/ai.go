package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type Resp struct {
	Result string `json:"result"`
}

func getAccessToken() string {
	// 使用APIKey和SecretKey获取access_token，替换下面的应用APIKey和应用SecretKey
	key := "b6xsGr4XQo0NgKmgBWUpXGTP"
	secret := "5ksZiuOEU2rbdkZNm3gKdgUsvH9zX3U6"
	url := fmt.Sprintf("https://aip.baidubce.com/oauth/2.0/token?grant_type=client_credentials&client_id=%s&client_secret=%s", key, secret)
	//url = fmt.Sprintf("https://aip.baidubce.com/oauth/2.0/token?client_id=%s&client_secret=%s&grant_type=client_credentials", key, secret)

	// 发送POST请求获取access_token
	resp, err := http.Post(url, "application/json", bytes.NewBuffer([]byte{}))
	if err != nil {
		fmt.Println("发送请求失败:", err)
		return ""
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取响应体失败:", err)
		return ""
	}

	// 解析JSON响应
	var data = make(map[string]interface{})
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("解析JSON失败:", err)
		return ""
	}

	// 获取access_token
	accessToken, ok := data["access_token"].(string)
	if !ok {
		fmt.Println("获取access_token失败")
		return ""
	}

	return accessToken
}

func chat(payload map[string]interface{}) {
	// 获取access_token
	accessToken := getAccessToken()
	if accessToken == "" {
		log.Fatalln("fail to get accessToken")
		return
	}

	// 构造请求URL
	url := "https://aip.baidubce.com/rpc/2.0/ai_custom/v1/wenxinworkshop/chat/ernie-4.0-8k-latest?access_token=" + accessToken

	// 构造请求体
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("构造请求体失败:", err)
		return
	}

	// 发送POST请求
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		fmt.Println("发送请求失败:", err)
		return
	}
	defer resp.Body.Close()

	// 逐行读取响应体内容
	reader := resp.Body
	buf := make([]byte, 1024)
	for {
		n, err := reader.Read(buf)
		if err != nil {
			if err == io.EOF {
				fmt.Println(string(buf[:n]))
				break
			}
			fmt.Println("读取响应体失败:", err)
			return
		}
		// 处理读取到的数据
		processData(buf[:n])
	}
}

func main() {
	payload := map[string]interface{}{
		"messages": []map[string]string{
			{
				"role":         "user",
				"content":      "你可以生成图片吗？",
				"enable_trace": "true",
			},
		},
		"stream": true,
	}
	chat(payload)
}

func processData(data []byte) {
	var resp Resp
	d := string(data)
	dl := strings.Split(d, "data:")
	b := []byte(dl[1])
	err := json.Unmarshal(b, &resp)
	if err != nil {
		log.Fatalln(err)
	}

	//fmt.Println(resp.Result)

	loop := []int32(resp.Result)

	for i := 0; i < len(loop); i++ {
		fmt.Println(string(loop[i]))
	}

}
