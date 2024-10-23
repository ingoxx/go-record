package main

import (
	"context"
	"fmt"
	pb "github.com/Lxb921006/go-record/grpc/streamrpc/streamrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"os"
	"path/filepath"
)

func main() {
	conn, err := grpc.NewClient(":12306", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	c := pb.NewMyServiceClient(conn)
	stream, err := c.MyMethod(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	file := "D:\\工作\\公司\\502打印机Brother MFC-8535DN驱动.EXE"
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
