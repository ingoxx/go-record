package main

import (
	"context"
	"fmt"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
)

// 初始化 Docker 客户端
func newDockerClient() (*client.Client, error) {
	return client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
}

// 初始化 Swarm 集群（仅适用于 Manager 节点）
func initSwarm(cli *client.Client, advertiseAddr string) {
	ctx := context.Background()

	// Swarm 初始化请求
	req := swarm.InitRequest{
		ListenAddr:    "0.0.0.0:2377",
		AdvertiseAddr: advertiseAddr,
		Spec: swarm.Spec{
			Annotations: swarm.Annotations{Name: "my-swarm-cluster"},
		},
	}

	nodeID, err := cli.SwarmInit(ctx, req)
	if err != nil {
		log.Fatalf("Failed to initialize Swarm: %v", err)
	}

	info, err := cli.Info(ctx)
	fmt.Println("Swarm initialized with Node ID:", nodeID, info.Swarm.Cluster.ID)
}

// 获取 Swarm 加入令牌
func getJoinTokens(cli *client.Client) (string, string) {
	ctx := context.Background()
	info, err := cli.SwarmInspect(ctx)
	if err != nil {
		log.Fatalf("Failed to get Swarm info: %v", err)
	}
	return info.JoinTokens.Worker, info.JoinTokens.Manager
}

// 添加节点到 Swarm（work 节点）
func joinSwarm(cli *client.Client, managerIP, token, advertiseAddr string) {
	ctx := context.Background()

	req := swarm.JoinRequest{
		ListenAddr:    "0.0.0.0:2377",
		AdvertiseAddr: advertiseAddr,
		RemoteAddrs:   []string{managerIP + ":2377"},
		JoinToken:     token,
	}

	err := cli.SwarmJoin(ctx, req)
	if err != nil {
		log.Fatalf("Failed to join Swarm: %v", err)
	}
	fmt.Println("Node successfully joined the Swarm")
}

// 获取 Swarm 节点列表（仅适用于 Manager 节点）
func listSwarmNodes(cli *client.Client) {
	ctx := context.Background()
	nodes, err := cli.NodeList(ctx, types.NodeListOptions{})
	if err != nil {
		log.Fatalf("Failed to list Swarm nodes: %v", err)
	}

	fmt.Println("Swarm Nodes:")
	for _, node := range nodes {
		fmt.Printf("ID: %s, Role: %s, Availability: %s, Status: %s\n",
			node.ID, node.Spec.Role, node.Spec.Availability, node.Status.State)
	}
}

// 删除 Swarm 节点（仅适用于 Manager 节点）
func removeNode(cli *client.Client, nodeID string) {
	ctx := context.Background()
	err := cli.NodeRemove(ctx, nodeID, types.NodeRemoveOptions{Force: true})
	if err != nil {
		log.Fatalf("Failed to remove node: %v", err)
	}
	fmt.Println("Node removed:", nodeID)
}

// 让当前节点退出 Swarm（Manager/work 节点）
func leaveSwarm(cli *client.Client, force bool) {
	ctx := context.Background()
	err := cli.SwarmLeave(ctx, force)
	if err != nil {
		log.Fatalf("Failed to leave Swarm: %v", err)
	}
	fmt.Println("Node left the Swarm")
}

// 关闭 Swarm（仅适用于 Manager 节点）
func removeSwarm(cli *client.Client) {
	ctx := context.Background()
	err := cli.SwarmLeave(ctx, true)
	if err != nil {
		log.Fatalf("Failed to remove Swarm: %v", err)
	}
	fmt.Println("Swarm removed")
}

// 主函数
func main() {
	cli, err := newDockerClient()
	if err != nil {
		log.Fatalf("Failed to create Docker client: %v", err)
	}
	defer cli.Close()

	managerIP := "192.168.1.100"
	workerIP := "192.168.1.101"

	// **步骤 1：初始化 Swarm**
	fmt.Println("Initializing Swarm...")
	initSwarm(cli, managerIP)

	// **步骤 2：获取 Swarm 加入令牌**
	workerToken, _ := getJoinTokens(cli)
	fmt.Println("Worker Join Token:", workerToken)

	// **步骤 3：在 worker 节点上加入 Swarm**
	workerClient, err := newDockerClient()
	if err != nil {
		log.Fatalf("Failed to create Docker client for worker: %v", err)
	}
	defer workerClient.Close()

	fmt.Println("Worker joining Swarm...")
	joinSwarm(workerClient, managerIP, workerToken, workerIP)

	// **步骤 4：列出 Swarm 节点**
	fmt.Println("Listing Swarm nodes...")
	listSwarmNodes(cli)

	// **步骤 5：删除 worker 节点**
	fmt.Println("Removing worker node...")
	nodes, err := cli.NodeList(context.Background(), types.NodeListOptions{})
	if err != nil {
		log.Fatalf("Failed to list nodes: %v", err)
	}

	for _, node := range nodes {
		if node.Spec.Role == swarm.NodeRoleWorker {
			removeNode(cli, node.ID)
		}
	}

	// **步骤 6：删除 Swarm**
	fmt.Println("Removing Swarm...")
	removeSwarm(cli)
}
