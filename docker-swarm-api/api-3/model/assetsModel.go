package model

import "time"

type AssetsModel struct {
	ID          uint         `json:"id" gorm:"primaryKey"`
	Ip          string       `json:"ip" gorm:"not null;unique"`
	NodeType    uint         `json:"node_type" gorm:"default:3;comment:1-master节点, 2-work节点, 3-未知节点类型"`
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
	ClusterID   *uint        `json:"cluster_id" gorm:"index;onDelete:SET NULL;default:NULL"`
	Cluster     ClusterModel `json:"cluster" gorm:"constraint:OnDelete:SET NULL;"`
	NodeStatus  uint         `json:"node_status" gorm:"default:300;comment:100-节点异常,200-节点正常,100-未知状态"`
}
