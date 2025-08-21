package ddw

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ingoxx/go-record/http/wx/pkg/config"
	"io"
	"log"
	"net/http"
	"time"
)

type DDWarn struct {
	msg string
}

func NewDDWarn(msg string) DDWarn {
	return DDWarn{
		msg: msg,
	}
}

func (dd DDWarn) paramReady() map[string]interface{} {
	return map[string]interface{}{
		"msgtype": "text",
		"text": map[string]interface{}{
			"content": fmt.Sprintf("小程序动态告警\n时间：%s\n内容: %s", time.Now().Format("2006-01-02 15:04:05"), dd.msg),
		},
		"at": map[string]interface{}{
			"atMobiles": "",
			"isAtAll":   false,
		},
	}
}

func (dd DDWarn) sendReq() error {
	b, err := json.Marshal(dd.paramReady())
	if err != nil {
		log.Printf("【ERROR】 请求钉钉接口时，序列化数据失败，失败信息：%s\n", err.Error())
		return err
	}

	resp, err := http.Post(config.AliWebHook, "application/json", bytes.NewBuffer(b))
	if err != nil {
		log.Printf("【ERROR】 请求钉钉接口失败，失败信息：%s\n", err.Error())
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("【ERROR】 请求钉钉接口信息失败，状态码：%d，失败信息：%s\n", resp.StatusCode, err.Error())
		return err
	}

	if _, err = io.ReadAll(resp.Body); err != nil {
		log.Printf("【ERROR】 读取钉钉接口信息失败，失败信息：%s\n", err.Error())
		return err
	}

	return nil
}

func (dd DDWarn) Send() error {
	return dd.sendReq()
}
