package process

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"

	Common "github.com/ingoxx/Golang-practise/net/practice01/common"
	"github.com/ingoxx/Golang-practise/net/practice01/model"
)

type FileProcess struct {
}

func (f *FileProcess) SendFileToUser(mes *model.Message) (err error) {
	fr := model.FileMessage{}
	err = json.Unmarshal([]byte(mes.Data), &fr)
	if err != nil {
		err = errors.New("反序列化SendFileToUser失败")
		return
	}

	up, k := UserMgr.OnLineUser[fr.UserId]
	if k {
		f.SendFileToOnLineUser(fr, up.Conn)
	}

	return
}

func (f *FileProcess) SendFileToOnLineUser(fr model.FileMessage, conn net.Conn) {
	mainmessage := model.Message{}
	mainmessage.Type = model.FileMessageType

	data, err := json.Marshal(&fr)
	if err != nil {
		fmt.Println("序列化SendFileToOnLineUser失败")
		return
	}

	mainmessage.Data = string(data)

	data, err = json.Marshal(&mainmessage)
	if err != nil {
		fmt.Println("序列化SendFileToOnLineUser失败")
		return
	}

	t := &Common.TransData{
		Conn: conn,
	}

	err = t.SendMessage(data)
	if err != nil {
		fmt.Println("SendFileToOnLineUser发消息失败, err= ", err)
		return
	}

}
