package main

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	pb "github.com/ingoxx/go-record/grpc/streamrpc/streamrpc"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"
)

type server struct {
	pb.UnimplementedStreamRpcServiceServer
	pb.UnimplementedMyServiceServer
}

func (s *server) SayHelloWorld(req *pb.StreamRequest, stream pb.StreamRpcService_SayHelloWorldServer) (err error) {
	log.Println("rec>>> ", req.GetName())

	for range [10]struct{}{} {
		if err = stream.Send(&pb.StreamReply{Message: "aaa"}); err != nil {
			return
		}
		time.Sleep(time.Duration(1) * time.Second)
	}

	return
}

func (s *server) MyMethod(stream pb.MyService_MyMethodServer) (err error) {
	log.Println("rec data")
	var done = make(chan struct{})

	if err = s.ProcessMsg(stream, done); err != nil {
		log.Println(err)
	}

	return
}

func (s *server) ProcessMsg(stream pb.MyService_MyMethodServer, done chan struct{}) (err error) {
	log.Println("ProcessData data")
	var file string
	var chunks = make([][]byte, 1024)

	for {
		resp, recErr := stream.Recv()
		if recErr == io.EOF {
			break
		}

		if resp == nil {
			return errors.New("nil pointer")
		}

		if file == "" {
			path := "./"
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

	log.Println(filepath.Base(file), " recv ok, returning to md5 soon")

	m, _ := s.FileMd5(file)

	if err = stream.Send(&pb.MyMessage{Msg: []byte(m), Name: filepath.Base(file)}); err != nil {
		return
	}

	//done <- struct_copy{}{}

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
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 12236))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	pb.RegisterStreamRpcServiceServer(s, &server{})
	pb.RegisterMyServiceServer(s, &server{})

	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
