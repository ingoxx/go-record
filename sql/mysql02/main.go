package main

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Tabler interface {
	TableName() string
}

type SmbModel struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Name      string `gorm:"not null" json:"name"`
	ShareName string `gorm:"not null" json:"shareName"`
}

func (SmbModel) TableName() string {
	return "smbServe_smbusermodel"
}

func main() {
	DB, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       "root:123321@tcp(43.134.182.215:34306)/cmdb?charset=utf8mb4&parseTime=True&loc=Local", // DSN data source name
		DefaultStringSize:         256,                                                                                   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,                                                                                  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,                                                                                  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,                                                                                  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,
	}), &gorm.Config{})
	if err != nil {
		log.Print(err)
		return
	}

	db, err := DB.DB()
	if err != nil {
		log.Print(err)
		return
	}

	defer db.Close()

	var sss []SmbModel

	uid := []uint{35}

	DB.Find(&sss, uid)

	fmt.Println(sss)

	DB.Where("id IN ?", uid).Unscoped().Delete(&SmbModel{})

}
