package main

import "time"

type ClusterModel struct {
	ID          uint          `json:"id" gorm:"primaryKey"`
	ClusterCid  string        `json:"cluster_cid" gorm:"default:n21q22l9bxkf0hhi7d971hh9o;comment:docker info可以查询"`
	Name        string        `json:"name" gorm:"unique"`
	Region      string        `json:"region" gorm:"default:cn-sz"`
	Date        time.Time     `json:"date" gorm:"default:CURRENT_TIMESTAMP;nullable"`
	Status      uint          `json:"status" gorm:"default:100;comment:100-集群异常,200-集群正常"`
	Servers     []AssetsModel `json:"servers" gorm:"foreignKey:ClusterID"`
	ClusterType string        `json:"cluster_type" gorm:"default:1"`
}

// AssetsModel 服务器列表
type AssetsModel struct {
	ID          uint         `json:"id" gorm:"primaryKey"`
	Ip          string       `json:"ip" gorm:"not null;unique"`
	NodeType    uint         `json:"node_type" gorm:"default:3;comment:1-master节点, 2-node节点, 3-未知节点类型"`
	Project     string       `json:"project" gorm:"not null"`
	Status      uint         `json:"status" gorm:"default:200;comment:100-服务器异常,200-服务器正常"`
	Operator    string       `json:"operator" gorm:"default:lxb"`
	RamUsage    uint         `json:"ram_usage" gorm:"default:1"`
	DiskUsage   uint         `json:"disk_usage" gorm:"default:1"`
	CpuUsage    uint         `json:"cpu_usage"  gorm:"default:1"`
	Start       time.Time    `json:"start" gorm:"default:CURRENT_TIMESTAMP;nullable"`
	User        string       `json:"user" gorm:"default:root"`
	Password    string       `json:"-" gorm:"not null"`
	Key         string       `json:"-" gorm:"type:TEXT"`
	Port        uint         `json:"port" gorm:"default:22"`
	OSType      uint         `json:"os_type" gorm:"default:1;comment:1-ubuntu,2-centos,3-debian"`
	ConnectType uint         `json:"connect_type" gorm:"default:1;comment:1-密码登陆, 2-秘钥登陆"`
	ClusterID   *uint        `json:"cluster_id" gorm:"index"`
	Cluster     ClusterModel `json:"cluster" gorm:"constraint:OnDelete:SET NULL;"`
}
