package process

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	Common "github.com/ingoxx/Golang-practise/net/practice01/common"
	"github.com/ingoxx/Golang-practise/net/practice01/model"
)

type FileProcess struct {
}

func (f *FileProcess) SendFile(path string, uid int) (err error) {
	mes := model.Message{}
	fr := model.FileMessage{}
	fn, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("%v不存在", path)
		return
	}

	if fn.Size() > int64(10485760) || fn.IsDir() {
		err = fmt.Errorf("%v必须是一个文件且大小不能超过10m", path)
		return
	}

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	mes.Type = model.FileMessageType
	fr.FileContent = b
	fr.FileName = fn.Name()
	fr.UserId = uid

	data, _ := json.Marshal(&fr)
	mes.Data = string(data)

	data, _ = json.Marshal(&mes)

	t := &Common.TransData{
		Conn: CurUsers.Conn,
	}

	err = t.SendMessage(data)
	if err != nil {
		err = fmt.Errorf("发送%v失败", path)
		return
	}

	return
}

func (f *FileProcess) SaveFile(fn string, content []byte) {
	path := "C:/Users/Administrator/Desktop/client"
	file := filepath.Join(path, fn)
	err := ioutil.WriteFile(file, content, 077)
	if err != nil {
		fmt.Println("SaveFile文件写入失败, err=", err)
	}
}
