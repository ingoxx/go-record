package main

import (
	"context"
	"fmt"
	pb "github.com/ingoxx/Golang-practise/grpc/command/command"
	"google.golang.org/grpc"
	"log"
	"net"
)

type Server struct {
	pb.UnimplementedUpdateProcessServer
}

func (s *Server) DockerUpdate(cxt context.Context, req *pb.CmdRequest) (resp *pb.CmdReply, err error) {
	log.Printf("DockerUpdate recv message: %s\n", req.GetMessage())
	return
}

func (s *Server) JavaUpdate(cxt context.Context, req *pb.CmdRequest) (resp *pb.CmdReply, err error) {
	log.Printf("JavaUpdate recv message: %s\n", req.GetMessage())
	return
}

func (s *Server) DockerReload(cxt context.Context, req *pb.CmdRequest) (resp *pb.CmdReply, err error) {
	log.Printf("DockerReload recv message: %s\n", req.GetMessage())
	return
}

func (s *Server) JavaReload(cxt context.Context, req *pb.CmdRequest) (resp *pb.CmdReply, err error) {
	log.Printf("JavaReload recv message: %s\n", req.GetMessage())
	return
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 12036))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterUpdateProcessServer(s, &Server{})

	log.Printf("rpc server listening at: %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
