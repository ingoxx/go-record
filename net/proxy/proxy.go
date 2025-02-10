package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net"
	"os"
)

var (
	emptyError = errors.New("empty json data")
)

type ProxyBackend struct {
	Listen  []map[string]string `json:"listen"`
	Backend []map[string]string `json:"backend"`
}

func handleClient(clientConn net.Conn, server string, servers []map[string]string) {
	log.Println("rec connect from ", clientConn.RemoteAddr())

	var remoteAddr string

	for _, v1 := range servers {
		sn, ok := v1[server]
		if ok {
			remoteAddr = sn
			break
		}
	}

	remoteConn, err := net.Dial("tcp", remoteAddr)
	if err != nil {
		log.Println("backend connect err >>> ", err)
		return
	}

	go func() {
		defer remoteConn.Close()
		defer clientConn.Close()
		_, err = io.Copy(remoteConn, clientConn)
		if err != nil {
			log.Println("remoteConn to clientConn err >>> ", err)
		}
	}()

	_, err = io.Copy(clientConn, remoteConn)
	if err != nil {
		log.Println("clientConn to remoteConn err >>> ", err)
	}

	log.Printf("proxy to %s finished\n", remoteAddr)
}

func main() {
	log.Println("--------------v1.0.0--------------")

	var pB ProxyBackend
	var stop = make(chan struct{})

	file, err := os.ReadFile("servers.json")
	if err != nil {
		log.Fatalln(err)
	}

	if err = json.Unmarshal(file, &pB); err != nil {
		log.Fatalln(err)
	}

	if len(pB.Listen) == 0 || len(pB.Backend) == 0 {
		log.Fatalln(emptyError)
	}

	for _, servers := range pB.Listen {
		for server, addr := range servers {
			log.Printf("start listen: %s   %s\n", server, addr)
			listener, err := net.Listen("tcp", addr) // 本地监听端口
			if err != nil {
				panic(err)
			}

			defer listener.Close()

			go func(server string) {
				for {
					clientConn, err := listener.Accept()
					if err != nil {
						panic(err)
					}
					go handleClient(clientConn, server, pB.Backend)
				}
			}(server)
		}
	}

	<-stop
}
