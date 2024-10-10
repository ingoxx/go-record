package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8088")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	name := []string{"lxb", "lqm", "lyy", "llt", "lch"}

	for _, v := range name {
		msg := fmt.Sprintf("Hello, server %s!", v)
		err = sendMessage(conn, msg)
		if err != nil {
			fmt.Println("Error sending message:", err)
		}
		time.Sleep(time.Second)
	}

}

func sendMessage(conn net.Conn, msg string) error {
	// 发送消息长度
	msgLen := uint32(len(msg))
	err := binary.Write(conn, binary.LittleEndian, msgLen)
	if err != nil {
		return err
	}

	// 发送消息内容
	_, err = conn.Write([]byte(msg))
	return err
}
