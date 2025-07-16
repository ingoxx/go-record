package common

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"net"

	UserMessage "github.com/ingoxx/Golang-practise/net/practice01/model"
)

type TransData struct {
	Conn net.Conn
	Buf  [102400]byte //这里还是换成切片好一些,再每次创建TransData实例时,初始化Buf就好了Buf: make([]bute, 1024)
}

func (t *TransData) RecvMessage() (mes UserMessage.Message, err error) {
	//buf := make([]byte, 8096)
	//读取客户端发来的数据, 如果客户端一直没发消息会阻塞,会出现超时
	_, err = t.Conn.Read(t.Buf[:4]) //从t.Conn中读取buf[:4]长度的消息放到buf
	if err != nil {
		// fmt.Printf("客户端=%v已退出\n", con.RemoteAddr().String())
		// err = errors.New("第一次接收数据失败")
		return
	}

	//显示到终端
	fmt.Printf("发送者=%v, 发送的内容=%v\n", t.Conn.RemoteAddr().String(), t.Buf[:4])

	//第一步先把接受的消息体长度转成uint32
	dl := binary.LittleEndian.Uint32(t.Buf[:4])

	//根据dl读取消息本体放到buf
	//fmt.Println("t.Buf[:dl] 消息内容", t.Buf[:dl]) //当dl大于len([1024]byte)会报错
	n, err := t.Conn.Read(t.Buf[:dl])  //这里读取的消息本体会重新覆盖放到buf
	if uint32(n) != dl || err != nil { //这里校验接收到的消息长度跟消息本体长度是否相同
		err = errors.New("第二次接收消息失败")
		// fmt.Println("消息丢失=", err)
		return
	}

	//上面如果消息长度校验通过则获取消息本体
	err = json.Unmarshal(t.Buf[:dl], &mes)
	if err != nil {
		err = errors.New("反序列化失败")
		// fmt.Println("反序列化失败=", err)
		return
	}

	return
}

func (t *TransData) SendMessage(data []byte) (err error) {
	//先发送长度给对方
	dl := uint32(len(data))
	fmt.Println("dl send len=", dl)
	//将消息长度转成一个表示长度的切片
	binary.LittleEndian.PutUint32(t.Buf[:4], dl)
	//发送消息长度
	n, e3 := t.Conn.Write(t.Buf[:4])
	if e3 != nil || n != 4 {
		err = errors.New("发送消息长度失败")
		return
	}

	//发送消息本身
	n, e3 = t.Conn.Write(data)
	if e3 != nil || n != int(dl) {
		err = errors.New("发送消息本体失败")
		return
	}

	return
}
