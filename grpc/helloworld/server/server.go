package main

import (
	"context"
	"flag"
	"fmt"
	pb "github.com/Lxb921006/Golang-practise/grpc/helloworld/helloworld"
	"google.golang.org/grpc"
	"log"
	"net"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

type server struct {
	pb.UnimplementedTestGrpcHelloWorldServer
}

func (s *server) SayHelloWorld(ctx context.Context, req *pb.HelloRequest) (resp *pb.HelloReply, err error) {
	log.Printf("recv message: %s", req.GetName())

	for range [10]struct{}{} {
		data := "hello " + req.GetName()
		r := &pb.HelloReply{
			Message: data,
		}
		resp = r
	}

	return
}

func main() {

	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterTestGrpcHelloWorldServer(s, &server{})

	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}