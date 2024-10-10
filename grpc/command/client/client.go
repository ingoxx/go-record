package client

import (
	"google.golang.org/grpc"
)

type Client struct {
	Conn *grpc.ClientConn
}

func (c *Client) DockerUpdate() (err error) {

	return
}
