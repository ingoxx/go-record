package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"log"

	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
)

// Worker 加入 Swarm
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

// 让节点退出 Swarm 集群
func leaveSwarm(cli *client.Client, force bool) {
	ctx := context.Background()

	// 执行节点退出 Swarm
	err := cli.SwarmLeave(ctx, force)
	if err != nil {
		log.Fatalf("Failed to leave Swarm: %v", err)
	}
	fmt.Println("Node left the Swarm")
}

// 初始化 Swarm 集群（仅适用于 Manager 节点）
func initSwarm(cli *client.Client, advertiseAddr string) {
	ctx := context.Background()

	// Swarm 初始化请求
	req := swarm.InitRequest{
		ListenAddr:    "0.0.0.0:2377",
		AdvertiseAddr: advertiseAddr,
		Spec: swarm.Spec{
			Annotations: swarm.Annotations{Name: "default"},
		},
	}

	nodeID, err := cli.SwarmInit(ctx, req)
	if err != nil {
		log.Fatalf("Failed to initialize Swarm: %v", err)
		//Failed to initialize Swarm: Error response from daemon: swarm spec must be named "default"
	}

	info, err := cli.Info(ctx)
	fmt.Println("Swarm initialized with Node ID:", nodeID, info.Swarm.Cluster.ID)

	tokens, s := getJoinTokens(cli)
	fmt.Printf("worker token %s, master token %s\n", tokens, s)
}

func listSwarmNodes(cli *client.Client) {
	ctx := context.Background()
	nodes, err := cli.NodeList(ctx, types.NodeListOptions{})
	if err != nil {
		log.Fatalf("Failed to list Swarm nodes: %v", err)
	}

	fmt.Println("Swarm Nodes:")
	for _, node := range nodes {
		nodeIP := node.Status.Addr
		nodeRole := node.Spec.Role
		nodeState := node.Status.State
		fmt.Println("nodeIP >>> ", nodeIP, nodeRole, swarm.NodeState(nodeState), node.ManagerStatus.Leader)
		fmt.Printf("ID: %s, Role: %s, Availability: %s, Status: %s\n",
			node.ID, node.Spec.Role, node.Spec.Availability, node.Status.State)
	}
}

func getJoinTokens(cli *client.Client) (string, string) {
	ctx := context.Background()
	info, err := cli.SwarmInspect(ctx)
	if err != nil {
		log.Fatalf("Failed to get Swarm info: %v", err)
	}
	return info.JoinTokens.Worker, info.JoinTokens.Manager
}

func main() {
	//managerIP := "127.0.0.1"
	//workerIP := "127.0.0.1"
	//workerToken := "SWMTKN-1-15cvjmq348altm36phkul207rfma3xboymldu9qsgf7u65ny37-8qv62v9skhsftkx8dp4f5n7a9" // 这个 token 需要从 manager 获取
	//
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalf("Failed to create Docker client: %v", err)
	}
	//defer cli.Close()
	//
	//fmt.Println("Worker joining Swarm...")
	//joinSwarm(cli, managerIP, workerToken, workerIP)
	//leaveSwarm(cli, true)
	listSwarmNodes(cli)
	//initSwarm(cli, "127.0.0.1")
}
