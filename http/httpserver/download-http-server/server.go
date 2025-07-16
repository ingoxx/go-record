package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"
)

func main() {
	httpServer()
}

type Resp struct {
	Msg    string                   `json:"msg"`
	Status int                      `json:"status"`
	Date   string                   `json:"date"`
	Detail map[string]interface{}   `json:"detail,omitempty"`
	Data   []map[string]interface{} `json:"data,omitempty"`
	Br     []byte                   `json:"br,omitempty"`
}

func (r *Resp) M(msg string, code int) (b []byte) {
	r.Msg = msg
	r.Status = code
	r.Date = time.Now().Add(time.Hour * time.Duration(11)).Format("2006-01-02 15:04:05")
	b, _ = json.Marshal(r)

	return
}

func (r *Resp) K(resp *Resp) (b []byte) {
	r.Msg = resp.Msg
	r.Status = resp.Status
	r.Date = time.Now().Add(time.Hour * time.Duration(11)).Format("2006-01-02 15:04:05")
	r.Detail = resp.Detail
	b, _ = json.Marshal(r)

	return
}

func (r *Resp) R(writer http.ResponseWriter, request *http.Request) error {
	_, err := writer.Write(r.Br)
	if err != nil {
		return err
	}

	log.Printf("[%v]  %s  %s\n", r.Status, request.RemoteAddr, request.URL)

	return nil
}

func (r *Resp) H(writer http.ResponseWriter, request *http.Request, data map[string]interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		log.Println("H() json序列化失败, ", err)
		return
	}
	_, err = writer.Write(b)
	if err != nil {
		log.Println("H() Write失败, ", err)
		return
	}
	log.Printf("[%v]  %s  %s\n", data["status"], request.RemoteAddr, request.URL)
	return
}

func httpServer() {
	log.Println("http server :9092 listening...")

	mux := http.NewServeMux()
	mux.HandleFunc("/download", download)
	mux.HandleFunc("/upload", upload)
	mux.HandleFunc("/content", sendFileContent)
	mux.HandleFunc("/aws-cdn-refresh", awsCdnRefresh)
	mux.HandleFunc("/wx-data", wxGetData)
	mux.HandleFunc("/change-time", changeTime)

	listen := &http.Server{
		Addr:              ":9092",
		Handler:           mux,
		ReadHeaderTimeout: time.Duration(10) * time.Second,
	}

	if err := listen.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf(err.Error())
	}
}

func changeTime(response http.ResponseWriter, request *http.Request) {
	var resp Resp
	if request.Method != "GET" {
		b := resp.M("请求方法错误", 10003)
		if _, err := response.Write(b); err != nil {
			log.Println("fail to response, esg: ", err)
			return
		}

		return
	}

	p := request.URL.Query()
	date := p.Get("date")
	sign := p.Get("sign")

	if sign != "change" {
		b := resp.M("无效请求", 10004)
		if _, err := response.Write(b); err != nil {
			return
		}

		return
	}

	log.Printf("date = %s, sign = %s\n", date, sign)

	scriptPath := "/web/wwwroot/777brs.com/web/change_time.sh"
	output, err := exec.Command("/bin/bash", scriptPath, date).Output()
	if err != nil {
		b := resp.M(fmt.Sprintf("执行%s失败，失败信息: %s", scriptPath, err), 10005)
		if _, err := response.Write(b); err != nil {
			return
		}

		return
	}

	b := resp.M(fmt.Sprintf("执行%s成功，成功信息: %s", scriptPath, string(output)), 10000)
	if _, err := response.Write(b); err != nil {
		return
	}

}

func wxGetData(resp http.ResponseWriter, req *http.Request) {
	var data = []map[string]interface{}{
		{
			"id":   1,
			"name": "lxb",
			"age":  31,
		},
		{
			"id":   2,
			"name": "lqm",
			"age":  18,
		},
		{
			"id":   3,
			"name": "lyy",
			"age":  17,
		},
	}

	var r Resp
	if req.Method != "GET" {
		b := r.M("请求方法错误", 10003)
		resp.Write(b)
		return
	}

	log.Println(req.RequestURI)

	f := req.URL.Query()
	var num = f.Get("num")
	li, _ := strconv.Atoi(num)
	if li < len(data) {
		data = data[:li]
	}

	r = Resp{
		Msg:    "ok",
		Status: 10000,
		Data:   data,
	}
	b := r.K(&r)

	resp.Write(b)

}

func upload(writer http.ResponseWriter, request *http.Request) {
	var resp Resp
	if request.Method != "POST" {
		resp.H(writer, request, map[string]interface{}{
			"esg":    "无效请求",
			"status": 10001,
		})
		return
	}

	if err := request.ParseMultipartForm(32 << 20); err != nil {
		resp.H(writer, request, map[string]interface{}{
			"esg":    err.Error(),
			"status": 10002,
		})
		return
	}

	file, header, _ := request.FormFile("file")
	fmt.Println("header >>> ", header)

	// 获取额外的参数
	value := request.Form.Get("user")
	fileName := request.Form.Get("fileName")

	fmt.Println(value, fileName)

	saveDir := filepath.Join("/nas/th-db-bak", fileName)
	fc, _ := os.Create(saveDir)

	defer fc.Close()

	if _, err := io.Copy(fc, file); err != nil {
		resp.H(writer, request, map[string]interface{}{
			"esg":    err.Error(),
			"status": 10003,
		})
		return
	}

	b := resp.M("上传成功", 10000)
	writer.Write(b)

}

func download(writer http.ResponseWriter, request *http.Request) {
	var resp Resp
	if request.Method != "GET" {
		b := resp.M("请求方法错误", 10003)
		writer.Write(b)
		return
	}

	f := request.URL.Query()

	if f.Get("token") != "123db" {
		b := resp.M("校验失败", 10001)
		writer.Write(b)
		return
	}

	if f.Get("file") == "" {
		b := resp.M("file字段不能为空", 10001)
		writer.Write(b)
		return
	}

	if err := sendFileHandle(f.Get("file"), writer); err != nil {
		b := resp.M(err.Error(), 10002)
		writer.Write(b)
		return
	}

	b := resp.M("ok", 10000)
	writer.Write(b)

}

func sendFileHandle(file string, w http.ResponseWriter) (err error) {
	fp := filepath.Join("/nas/tempDB", file)
	f, err := os.Open(fp)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file))
	w.Header().Set("Content-Type", "application/octet-stream")
	//w.Header().Set("Content-Length", "123456789") // 设置文件大

	_, err = io.Copy(w, f)
	if err != nil {
		return
	}

	return
}

func sendFileContent(writer http.ResponseWriter, request *http.Request) {
	var resp Resp
	if request.Method != "GET" {
		b := resp.M("请求方法错误", 10003)
		writer.Write(b)
		return
	}

	f := request.URL.Query()
	file := filepath.Join("C:\\Users\\Administrator\\Desktop", f.Get("file"))
	_, err := os.Stat(file)
	if err != nil {
		b := resp.M(err.Error(), 10002)
		writer.Write(b)
		return
	}

	http.ServeFile(writer, request, file)

}

func awsCdnRefresh(writer http.ResponseWriter, request *http.Request) {
	//log.Printf("receive req: %s %s %s\n", request.RemoteAddr, request.URL, request.Method)
	var resp Resp
	var awsResp map[string]interface{}
	if request.Method != "GET" {
		resp.H(writer, request, map[string]interface{}{
			"esg":    "无效请求",
			"status": 10001,
		})
		return
	}

	f := request.URL.Query()
	path := f.Get("path")
	item := f.Get("item")
	sign := f.Get("sign")
	if path == "" || item == "" || sign == "" {
		resp.H(writer, request, map[string]interface{}{
			"esg":    "参数错误",
			"status": 10002,
		})
		return
	}

	if sign != "pyhdiaomao" {
		resp.H(writer, request, map[string]interface{}{
			"esg":    "签名错误",
			"status": 10003,
		})
		return
	}

	var cancel context.CancelFunc
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(50))
	defer cancel()

	out, err := exec.CommandContext(ctx, "bash", "/root/shellscript/aws_cdn_refresh.sh", item, path).Output()
	if err != nil {
		log.Println("exec err: ", err, string(out))
		resp.H(writer, request, map[string]interface{}{
			"esg":    string(out),
			"status": 10004,
		})
		return
	}

	if err = json.Unmarshal(out, &awsResp); err != nil {
		log.Println("unmarshal err: ", err, string(out))
		resp.H(writer, request, map[string]interface{}{
			"esg":    err.Error(),
			"status": 10005,
		})
		return
	}

	respKRM := &Resp{
		Msg:    fmt.Sprintf("%s cdn refresh succeed", item),
		Status: 10000,
		Detail: awsResp,
	}

	b := resp.K(respKRM)
	respKRM.Br = b

	if err := respKRM.R(writer, request); err != nil {
		log.Println("响应失败, 失败信息: ", err)
		return
	}
}
