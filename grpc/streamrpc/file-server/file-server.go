package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	pb "github.com/ingoxx/go-record/grpc/streamrpc/streamrpc"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
)

type server struct {
	pb.UnimplementedMyServiceServer
}

func (s *server) MyMethod(stream pb.MyService_MyMethodServer) (err error) {
	log.Println("rec data")

	if err = s.ProcessData(stream); err != nil {
		log.Println(err)
	}

	return
}

func (s *server) ProcessData(stream pb.MyService_MyMethodServer) (err error) {
	log.Println("process msg")
	//stream := <-s.work

	var file string
	var chunks [][]byte

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			log.Println("rec finished")
			break
		}

		if file == "" {
			path := "C:\\Users\\Administrator\\Desktop"
			file = filepath.Join(path, resp.GetName())
		}

		chunks = append(chunks, resp.Msg)
	}

	_, err = os.Stat(file)
	if err != nil {
		fw, err := os.Create(file)
		if err != nil {
			return err
		}
		defer fw.Close()

		for _, chunk := range chunks {
			_, err := fw.WriteString(string(chunk))
			if err != nil {
				return err
			}
		}
	}

	log.Println(file, " recv ok")

	m, _ := s.FileMd5(file)

	if err = stream.Send(&pb.MyMessage{Msg: []byte("md5"), Name: m}); err != nil {
		log.Println("send err ", err)
		return
	}

	return
}

func (s *server) FileMd5(file string) (m5 string, err error) {
	f, err := os.Open(file)
	if err != nil {
		return
	}

	defer f.Close()

	h := md5.New()
	if _, err = io.Copy(h, f); err != nil {
		return
	}

	m5 = hex.EncodeToString(h.Sum(nil))

	return
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 12306))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterMyServiceServer(s, &server{})

	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
