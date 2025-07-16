package main

import (
	"github.com/ingoxx/go-record/docker-swarm-api/api-3/check"
	"log"

	"github.com/docker/docker/client"
)

func main() {
	// 创建 Docker 客户端
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalf("Failed to create Docker client: %v", err)
	}
	defer cli.Close()

	// 运行 Swarm 健康检查
	check.Health(cli)
}
