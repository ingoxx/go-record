package main

import (
	"context"
	"flag"
	"fmt"
	pb "github.com/ingoxx/go-record/grpc/streamrpc/streamrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"os"
	"path/filepath"
)

func main() {
	addr := flag.String("addr", "", "rpc服务端地址端口") // 定义一个字符串类型的标志
	file := flag.String("file", "30", "文件路径")     // 定义一个整数类型的标志
	flag.Parse()
	if *addr == "" || *file == "" {
		_, err := fmt.Fprintln(os.Stderr, "Error: -addr跟-file都是必须参数，不能为空")
		if err != nil {
			return
		}
		os.Exit(1)
	}
	processData(*addr, *file)
}

func processData(addr, file string) {
	conn, err := grpc.NewClient(fmt.Sprintf("%s", addr), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	c := pb.NewMyServiceClient(conn)
	stream, err := c.MyMethod(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fr, err := os.Open(file)
	if err != nil {
		log.Fatalln(err)
	}

	var fb = make([]byte, 1024)
	for {
		n, err := fr.Read(fb)
		if err == io.EOF {
			break
		}
		if n == 0 {
			break
		}
		if err = stream.Send(&pb.MyMessage{Msg: fb[:n], Name: filepath.Base(file)}); err != nil {
			return
		}
	}

	stream.CloseSend()

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("file: %s, md5: %s", resp.Name, string(resp.Msg))
	}
}
