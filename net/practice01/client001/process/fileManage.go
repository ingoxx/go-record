package process

import (
	"encoding/json"
	"fmt"

	UserMessage "github.com/ingoxx/Golang-practise/net/practice01/model"
)

// 接收服务端推送的离线聊天记录
func RecvFile(mes *UserMessage.Message) {
	filedata := UserMessage.FileMessage{}
	err := json.Unmarshal([]byte(mes.Data), &filedata)
	if err != nil {
		fmt.Println("反序列化RecvFile失败,err= ", err)
		return
	}

	fp := &FileProcess{}
	fp.SaveFile(filedata.FileName, filedata.FileContent)
	fmt.Printf("用户id:%d, 发了个文件:%v", filedata.UserId, filedata.FileName)
}
