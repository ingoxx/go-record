package model

import (
	"github.com/ingoxx/go-record/docker-swarm-api/api-3/db"
	"time"
)

type ClusterModel struct {
	ID          uint          `json:"id" gorm:"primaryKey"`
	ClusterCid  string        `json:"cluster_cid" gorm:"default:n21q22l9bxkf0hhi7d971hh9o;comment:docker info可以查询"`
	Name        string        `json:"name" gorm:"unique"`
	Region      string        `json:"region" gorm:"default:cn-sz"`
	WorkToken   string        `json:"-" gorm:"null;comment:work节点的token"`
	MasterToken string        `json:"-" gorm:"null;comment:master节点token"`
	MasterIp    string        `json:"master_ip"  gorm:"default:1.1.1.1"`
	Date        time.Time     `json:"date" gorm:"default:CURRENT_TIMESTAMP;nullable"`
	Status      uint          `json:"status" gorm:"default:200;comment:100-集群异常,200-集群正常"`
	Servers     []AssetsModel `json:"servers" gorm:"foreignKey:ClusterID"`
	ClusterType string        `json:"cluster_type" gorm:"default:1"`
}

func (cm *ClusterModel) Update(id uint) error {
	if err := db.DB.Where("id = ?", id).Delete(cm).Error; err != nil {
		return err
	}

	return nil
}
