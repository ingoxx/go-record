package main

import (
	"context"
	"flag"
	pb "github.com/ingoxx/go-record/grpc/helloworld/helloworld"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	name = flag.String("name", "lxb", "Name to sayHelloWorld")
)

func main() {
	//flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	c := pb.NewTestGrpcHelloWorldClient(conn)

	r, err := c.SayHelloWorld(ctx, &pb.HelloRequest{Name: *name})
	if err != nil {
		log.Fatalf("could not send: %v", err)
	}

	// Contact the server and print out its response.
	log.Printf("recv: %s", r.GetMessage())
	time.Sleep(time.Second)

}
