package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/ingoxx/Golang-practise/redis/redisServer"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Resp struct {
	Result string `json:"result"`
}

type WxLogin struct {
	Openid string `json:"openid"`
}

type QuotaInfo struct {
	ChatGpt  int `json:"chatgpt"`
	Qw       int `json:"qw"`
	Gemini   int `json:"gemini"`
	Bd       int `json:"bd"`
	Invite   int `json:"invite"`
	Finished int `json:"finished"`
	Time     int `json:"time"`
	ErrCode  int `json:"errcode"`
}

type wsData struct {
	Model   string              `json:"model"`
	Title   string              `json:"title"`
	Context []map[string]string `json:"context"`
}

// 公共的错误
func errorData(data []rune, ws *websocket.Conn, messageType int) {
	for _, v := range data {
		err := ws.WriteMessage(messageType, []byte(string(v)))
		if err != nil {
			log.Println("Error during message writing:", err)
			break
		}
	}

}

// 百度
func getAccessToken() string {
	// 使用APIKey和SecretKey获取access_token，替换下面的应用APIKey和应用SecretKey
	key := "b6xsGr4XQo0NgKmgBWUpXGTP"
	secret := "5ksZiuOEU2rbdkZNm3gKdgUsvH9zX3U6"
	url := fmt.Sprintf("https://aip.baidubce.com/oauth/2.0/token?grant_type=client_credentials&client_id=%s&client_secret=%s", key, secret)

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
	var data map[string]interface{}
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

func chatBaiDu(payload map[string]interface{}, ws *websocket.Conn, messageType int, message []byte) {
	var wsData wsData
	json.Unmarshal(message, &wsData)

	// 获取access_token
	accessToken := getAccessToken()
	if accessToken == "" {
		return
	}

	// 构造请求URL
	un := "https://aip.baidubce.com/rpc/2.0/ai_custom/v1/wenxinworkshop/chat/ernie-4.0-8k-latest?access_token=" + accessToken

	payload = map[string]interface{}{
		"messages":     wsData.Context,
		"stream":       true,
		"enable_trace": true,
	}

	// 构造请求体
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("构造请求体失败:", err)
		return
	}

	// 发送POST请求
	resp, err := http.Post(un, "application/json", bytes.NewBuffer(payloadBytes))
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
				break
			}
			fmt.Println("读取响应体失败:", err)
			return
		}
		time.Sleep(time.Millisecond / 40)
		// 处理读取到的数据
		processBaiDuData(buf[:n], ws, messageType)
	}
}

func processBaiDuData(data []byte, ws *websocket.Conn, messageType int) {
	var resp Resp
	d := string(data)
	dl := strings.Split(d, "data:")

	b := []byte(dl[1])
	err := json.Unmarshal(b, &resp)
	if err != nil {
		log.Fatalln("Unmarshal failed >>> ", err)
	}

	loop := []rune(resp.Result)

	for _, v := range loop {
		err := ws.WriteMessage(messageType, []byte(string(v)))
		if err != nil {
			log.Println("Error during message writing:", err)
			break
		}
		time.Sleep(time.Millisecond / 40)
	}

}

func bDWebsocketHandler(w http.ResponseWriter, r *http.Request) {
	var msg = make(map[string]interface{})
	var q = r.URL.Query()
	log.Println("bd openid >>> ", q.Get("openid"))
	if q.Get("openid") == "" {
		msg["esg"] = "openid不能为空"
		p, _ := json.Marshal(&msg)
		_, _ = w.Write(p)
		return
	}

	var rds = redisServer.NewRds(q.Get("openid"))
	if err := rds.CheckOpenId(); err != nil {
		msg["esg"] = err.Error()
		p, _ := json.Marshal(&msg)
		_, _ = w.Write(p)
		return
	}

	conn, err := upGrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("Error during connection upgradation:", err)
		return
	}
	defer conn.Close()

	messageType, message, err := conn.ReadMessage()
	if err != nil {
		log.Println("Error during message reading:", err)
		return
	}
	log.Printf("bd Received: %s", message)

	if err := rds.UpdateQuota("bd"); err != nil {
		bt := []rune(err.Error())
		log.Printf("bd ai limit >>> %s\n", err.Error())
		errorData(bt, conn, messageType)
		return
	}

	var payload map[string]interface{}
	chatBaiDu(payload, conn, messageType, message)
}

// 谷歌
func geminiWebsocketHandler(w http.ResponseWriter, r *http.Request) {
	var msg = make(map[string]interface{})
	var q = r.URL.Query()
	var openid = q.Get("openid")
	log.Println("gemini openid >>> ", openid)
	if q.Get("openid") == "" {
		msg["esg"] = "openid不能为空"
		p, _ := json.Marshal(&msg)
		_, _ = w.Write(p)
		return
	}

	var rds = redisServer.NewRds(openid)
	if err := rds.CheckOpenId(); err != nil {
		msg["esg"] = err.Error()
		p, _ := json.Marshal(&msg)
		_, _ = w.Write(p)
		return
	}

	conn, err := upGrader.Upgrade(w, r, nil)

	if err != nil {
		log.Print("Error during connection upgradation:", err)
		return
	}

	defer conn.Close()

	messageType, message, err := conn.ReadMessage()
	if err != nil {
		log.Println("Error during message reading:", err)
		return
	}

	log.Printf("gemini Received: %s", message)

	if err := rds.UpdateQuota("gemini"); err != nil {
		bt := []rune(err.Error())
		log.Printf("gemini ai limit >>> %s\n", err.Error())
		errorData(bt, conn, messageType)
		return
	}

	chatGemini(conn, messageType, message, openid)
}

func chatGemini(ws *websocket.Conn, messageType int, message []byte, openid string) {
	var wsData wsData
	_ = json.Unmarshal(message, &wsData)

	socketUrl := fmt.Sprintf("ws://127.0.0.1:10086/gemini/%s/", openid)
	conn, _, err := websocket.DefaultDialer.Dial(socketUrl, nil)
	if err != nil {
		log.Println("Error connecting to Websocket Server:", err)
		return
	}

	defer conn.Close()

	err = conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		log.Println("Error during writing to websocket:", err)
		return
	}

	for {
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error in receive:", err)
			return
		}
		processGeminiData(msg, ws, messageType)
	}

}

func processGeminiData(data []byte, ws *websocket.Conn, messageType int) {
	var resp Resp

	err := json.Unmarshal(data, &resp)
	if err != nil {
		log.Println("Unmarshal failed >>> ", err)
		return
	}

	loop := []rune(resp.Result)

	for _, v := range loop {
		err := ws.WriteMessage(messageType, []byte(string(v)))
		if err != nil {
			log.Println("Error during message writing:", err)
			break
		}
		//time.Sleep(time.Millisecond / 70)
	}
}

// 通义千问
func qWWebsocketHandler(w http.ResponseWriter, r *http.Request) {
	var msg = make(map[string]interface{})
	var q = r.URL.Query()
	var openid = q.Get("openid")
	log.Println("qw openid >>> ", openid)
	if q.Get("openid") == "" {
		msg["esg"] = "openid不能为空"
		p, _ := json.Marshal(&msg)
		_, _ = w.Write(p)
		return
	}

	var rds = redisServer.NewRds(openid)
	if err := rds.CheckOpenId(); err != nil {
		msg["esg"] = err.Error()
		p, _ := json.Marshal(&msg)
		_, _ = w.Write(p)
		return
	}

	if err := rds.UpdateQuota("qw"); err != nil {
		msg["esg"] = err.Error()
		p, _ := json.Marshal(&msg)
		_, _ = w.Write(p)
		return
	}

	conn, err := upGrader.Upgrade(w, r, nil)

	if err != nil {
		log.Print("Error during connection upgradation:", err)
		return
	}

	defer conn.Close()

	messageType, message, err := conn.ReadMessage()
	if err != nil {
		log.Println("Error during message reading:", err)
		return
	}

	log.Printf("qw Received: %s", message)

	if err := rds.UpdateQuota("qw"); err != nil {
		bt := []rune(err.Error())
		log.Printf("chatgpt ai limit >>> %s\n", err.Error())
		errorData(bt, conn, messageType)
		return
	}

	chatQw(conn, messageType, message, openid)
}

func processQwData(data []byte, ws *websocket.Conn, messageType int) {
	var resp Resp

	err := json.Unmarshal(data, &resp)
	if err != nil {
		log.Fatalln("Unmarshal failed >>> ", err)
	}

	loop := []rune(resp.Result)

	for _, v := range loop {
		err := ws.WriteMessage(messageType, []byte(string(v)))
		if err != nil {
			log.Println("Error during message writing:", err)
			break
		}
		time.Sleep(time.Millisecond / 40)
	}
}

func chatQw(ws *websocket.Conn, messageType int, message []byte, openid string) {
	var wsData wsData
	_ = json.Unmarshal(message, &wsData)

	dial := websocket.Dialer{
		HandshakeTimeout: 300 * time.Second,
	}

	socketUrl := fmt.Sprintf("ws://127.0.0.1:10086/qw/%s/", openid)

	conn, _, err := dial.Dial(socketUrl, nil)
	if err != nil {
		log.Fatal("Error connecting to Websocket Server:", err)
	}

	defer conn.Close()

	err = conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		log.Fatal("Error during writing to websocket:", err)
	}

	for {
		//接收
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error in receive:", err)
			return
		}
		processQwData(msg, ws, messageType)
	}

}

// chatGpt
func chatGptWebsocketHandler(w http.ResponseWriter, r *http.Request) {
	var msg = make(map[string]interface{})
	var q = r.URL.Query()
	var openid = q.Get("openid")
	log.Println("chatgpt openid >>> ", openid)
	if q.Get("openid") == "" {
		msg["esg"] = "openid不能为空"
		p, _ := json.Marshal(&msg)
		_, _ = w.Write(p)
		return
	}

	var rds = redisServer.NewRds(openid)
	if err := rds.CheckOpenId(); err != nil {
		msg["esg"] = err.Error()
		p, _ := json.Marshal(&msg)
		_, _ = w.Write(p)
		return
	}

	if err := rds.UpdateQuota("chatgpt"); err != nil {
		msg["esg"] = err.Error()
		p, _ := json.Marshal(&msg)
		_, _ = w.Write(p)
		return
	}

	conn, err := upGrader.Upgrade(w, r, nil)

	if err != nil {
		log.Print("Error during connection upgradation:", err)
		return
	}

	defer conn.Close()

	messageType, message, err := conn.ReadMessage()
	if err != nil {
		log.Println("Error during message reading:", err)
		return
	}

	log.Printf("chatgpt Received: %s", message)

	if err := rds.UpdateQuota("chatgpt"); err != nil {
		bt := []rune(err.Error())
		log.Printf("chatgpt ai limit >>> %s\n", err.Error())
		errorData(bt, conn, messageType)
		return
	}

	chatGpt(conn, messageType, message, openid)
}

func processChatGptData(data []byte, ws *websocket.Conn, messageType int) {
	var resp Resp

	err := json.Unmarshal(data, &resp)
	if err != nil {
		log.Fatalln("Unmarshal failed >>> ", err)
	}

	loop := []rune(resp.Result)

	for _, v := range loop {
		err := ws.WriteMessage(messageType, []byte(string(v)))
		if err != nil {
			log.Println("Error during message writing:", err)
			break
		}
		time.Sleep(time.Millisecond / 40)
	}
}

func chatGpt(ws *websocket.Conn, messageType int, message []byte, openid string) {
	var wsData wsData
	_ = json.Unmarshal(message, &wsData)

	dial := websocket.Dialer{
		HandshakeTimeout: 300 * time.Second,
	}

	socketUrl := fmt.Sprintf("ws://127.0.0.1:10086/chatgpt/%s/", openid)

	conn, _, err := dial.Dial(socketUrl, nil)
	if err != nil {
		log.Fatal("Error connecting to Websocket Server:", err)
	}

	defer conn.Close()

	err = conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		log.Fatal("Error during writing to websocket:", err)
	}

	for {
		//接收
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error in receive:", err)
			return
		}
		processChatGptData(msg, ws, messageType)
	}

}

// 微信登录
func getWxOpenId(resp http.ResponseWriter, req *http.Request) {
	log.Println(req.Method, " getWxOpenId")
	var msg = make(map[string]interface{})
	var codeData map[string]interface{}
	if req.Method != "POST" {
		msg["esg"] = "无效请求"
		msg["errcode"] = 10001
		p, _ := json.Marshal(&msg)
		_, _ = resp.Write(p)
		return
	}

	bd, err := io.ReadAll(req.Body)
	if err != nil {
		msg["esg"] = "解析请求失败"
		msg["errcode"] = 10002
		p, _ := json.Marshal(&msg)
		_, _ = resp.Write(p)
		return
	}

	_ = json.Unmarshal(bd, &codeData)

	v := url.Values{}
	v.Add("appid", "wxbb1377eff3149db4")
	v.Add("secret", "6bae191e5e03aa4cf5731478ab513624")
	v.Add("js_code", codeData["code"].(string))
	v.Add("grant_type", "authorization_code")

	urlName := "https://api.weixin.qq.com/sns/jscode2session?" + v.Encode()

	r, err := http.Get(urlName)
	defer r.Body.Close()

	if err != nil {
		msg["esg"] = err.Error()
		msg["errcode"] = 10002
		p, _ := json.Marshal(&msg)
		resp.Write(p)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		msg["esg"] = err.Error()
		msg["errcode"] = 10003
		p, _ := json.Marshal(&msg)
		resp.Write(p)
		return
	}

	if r.StatusCode != 200 {
		resp.Write(body)
		return
	}

	var data WxLogin
	json.Unmarshal(body, &data)
	var openid = data.Openid
	var rds = redisServer.NewRds(openid)
	if rds.Err != nil {
		msg["esg"] = err.Error()
		msg["errcode"] = 10004
		p, _ := json.Marshal(&msg)
		resp.Write(p)
		return
	}

	// 初始化一些参数
	err = rds.Set()
	err = rds.SetQuota()
	err = rds.SetInvite()

	b, _ := json.Marshal(&data)

	resp.Write(b)

}

// 额度查询
func getQuota(resp http.ResponseWriter, req *http.Request) {
	log.Println(req.Method, " getQuota")
	var msg = make(map[string]interface{})
	if req.Method != "GET" {
		msg["esg"] = "无效请求"
		msg["errcode"] = 10001
		p, _ := json.Marshal(&msg)
		_, err := resp.Write(p)
		log.Println("getQuota Method err >>> ", err)
		return
	}
	var q = req.URL.Query()
	var openid = q.Get("openid")

	if openid == "" {
		msg["esg"] = "openid不能为空"
		msg["errcode"] = 10002
		p, _ := json.Marshal(&msg)
		_, _ = resp.Write(p)
		return
	}

	var rds = redisServer.NewRds(openid)
	data, err := rds.GetQuota()
	if err != nil {
		msg["esg"] = err.Error()
		msg["errcode"] = 10003
		p, _ := json.Marshal(&msg)
		_, _ = resp.Write(p)
		return
	}

	log.Println("data >>> ", data)

	var respData = QuotaInfo{
		ChatGpt:  data["chatgpt"],
		Qw:       data["qw"],
		Gemini:   data["gemini"],
		Bd:       data["bd"],
		Invite:   data["invite"],
		Finished: data["finished"],
		Time:     data["time"],
		ErrCode:  10000,
	}

	b, _ := json.Marshal(&respData)
	_, _ = resp.Write(b)

}

// 邀请
func invite(resp http.ResponseWriter, req *http.Request) {
	log.Println(req.Method, " invite")
	var msg = make(map[string]interface{})
	if req.Method != "GET" {
		msg["esg"] = "无效请求"
		msg["errcode"] = 10001
		p, _ := json.Marshal(&msg)
		_, err := resp.Write(p)
		log.Println("invite Method err >>> ", err)
		return
	}

	var q = req.URL.Query()
	var openid = q.Get("openid")
	var userid = q.Get("userid")

	if openid == "" || userid == "" {
		msg["esg"] = "openid或userid不能为空"
		msg["errcode"] = 10002
		p, _ := json.Marshal(&msg)
		_, err := resp.Write(p)
		log.Println("invite 1 resp err >>> ", err)
		return
	}

	var rds = redisServer.NewRds(openid)
	in, err := rds.Invite(userid)
	if err != nil {
		log.Println("invite 2 resp err >>> ", err)
		msg["esg"] = err.Error()
		msg["errcode"] = 10003
		p, _ := json.Marshal(&msg)
		_, _ = resp.Write(p)
		return
	}

	msg["esg"] = "ok"
	msg["invite"] = in
	msg["errcode"] = 10000
	p, _ := json.Marshal(&msg)
	_, _ = resp.Write(p)
}

// 上传文件
func upload(resp http.ResponseWriter, req *http.Request) {
	log.Println(req.Method, " upload")
	var msg = struct {
		Esg     string `json:"esg"`
		ErrCode int    `json:"errcode"`
	}{}

	if req.Method != "POST" {
		msg.Esg = "无效请求"
		msg.ErrCode = 10001
		p, _ := json.Marshal(&msg)
		_, _ = resp.Write(p)
		return
	}

	form := req.ParseMultipartForm(32 << 20)
	if form != nil {
		msg.Esg = "文件大小超出最大限制32M"
		msg.ErrCode = 10002
		p, _ := json.Marshal(&msg)
		_, _ = resp.Write(p)
		return
	}

	var q = req.URL.Query()
	var openid = q.Get("openid")
	if openid == "" {
		msg.Esg = "openid不能为空"
		msg.ErrCode = 10003
		p, _ := json.Marshal(&msg)
		_, _ = resp.Write(p)
		return
	}

	var rds = redisServer.NewRds(openid)
	if err := rds.CheckOpenId(); err != nil {
		msg.Esg = err.Error()
		msg.ErrCode = 10004
		p, _ := json.Marshal(&msg)
		_, _ = resp.Write(p)
		return
	}

	file, _, _ := req.FormFile("file")

	// 获取额外的参数
	fileName := req.Form.Get("fileName")

	saveDir := filepath.Join("/data/chat/web", fileName)
	fc, _ := os.Create(saveDir)

	defer fc.Close()

	_, err := io.Copy(fc, file)
	if err != nil {
		msg.Esg = err.Error()
		msg.ErrCode = 10005
		p, _ := json.Marshal(&msg)
		_, _ = resp.Write(p)
		return
	}

	msg.Esg = "上传成功"
	msg.ErrCode = 10000
	p, err := json.Marshal(&msg)
	_, _ = resp.Write(p)

}

func main() {
	log.Println("listening ::10087")
	http.HandleFunc("/bd", bDWebsocketHandler)
	http.HandleFunc("/gemini", geminiWebsocketHandler)
	http.HandleFunc("/chatgpt", chatGptWebsocketHandler)
	http.HandleFunc("/qw", qWWebsocketHandler)
	http.HandleFunc("/wx-login", getWxOpenId)
	http.HandleFunc("/get-quota", getQuota)
	http.HandleFunc("/invite", invite)
	http.HandleFunc("/upload", upload)
	log.Fatal(http.ListenAndServe("10.0.0.15:10087", nil))
}
