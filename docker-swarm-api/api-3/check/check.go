package check

import (
	"context"
	"fmt"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// Health 检测 Swarm 集群健康状态
func Health(cli *client.Client) {
	ctx := context.Background()

	// 获取 Swarm 集群中的所有节点
	nodes, err := cli.NodeList(ctx, types.NodeListOptions{})
	if err != nil {
		log.Fatalf("Failed to get node list: %v", err)
	}

	// 统计 Manager 和 Worker 状态
	healthyManagers := 0
	totalManagers := 0
	unhealthyNodes := 0

	for _, node := range nodes {
		role := node.Spec.Role
		status := node.Status.State
		availability := node.Spec.Availability
		managerStatus := "Worker"

		// 如果是 Manager，则检查其状态
		if node.ManagerStatus != nil {
			managerStatus = "Manager"
			totalManagers++
			if node.ManagerStatus.Reachability == "Reachable" {
				healthyManagers++
			}
		}

		// 输出节点信息
		fmt.Printf("Node: %s | Role: %s | Status: %s | Availability: %s | Manager: %s\n",
			node.Description.Hostname, role, status, availability, managerStatus)

		// 如果节点状态异常，增加计数
		if status != "Ready" {
			unhealthyNodes++
		}
	}

	// 评估 Swarm 健康状态
	if unhealthyNodes > 0 {
		fmt.Println("❌ 警告：集群中有", unhealthyNodes, "个异常节点")
	} else {
		fmt.Println("✅ 所有节点状态正常")
	}

	// 检查 Manager 是否足够
	if totalManagers > 0 && healthyManagers == 0 {
		fmt.Println("❌ 警告：所有 Manager 都不可用，Swarm 可能无法运作")
	} else if totalManagers > 0 {
		fmt.Println("✅ Manager 状态正常（", healthyManagers, "/", totalManagers, "）")
	}
}
