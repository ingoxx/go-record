package main

import (
	"encoding/binary"
	"io"
	"log"
	"net"
)

func main() {
	log.Println("server listening on 8088...")

	listener, err := net.Listen("tcp", ":8088")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	for {
		// 读取消息长度
		var msgLen uint32
		err := binary.Read(conn, binary.LittleEndian, &msgLen)
		if err == io.EOF {
			log.Println(conn.RemoteAddr().String(), " client already disconnect")
			return
		}

		if err != nil {
			log.Println(conn.RemoteAddr().String(), " Error reading message length: ", err)
			return
		}

		// 读取消息内容
		buf := make([]byte, msgLen)
		_, err = io.ReadFull(conn, buf)
		if err != nil {
			log.Println(conn.RemoteAddr().String(), " Error reading message:", err)
			return
		}

		log.Printf("Received message from %s: %s", conn.RemoteAddr().String(), string(buf))
	}
}
