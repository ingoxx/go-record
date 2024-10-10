package main

import (
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type CronsCrontabs struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	OperateUser string    `gorm:"not null" json:"operator"`
	Mission     string    `gorm:"not null" json:"cron"`
	StartTime   time.Time `json:"st"`
	EndTime     time.Time `json:"et"`
	Status      uint      `gorm:"default:100" json:"status"`
	Project     string    `gorm:"not null" json:"project"`
}

type Cron struct {
}

func main() {
	// sqlDB, err := sql.Open("mysql", "root:123321@tcp(43.138.184.202:34306)/cmdb?charset=utf8mb4&parseTime=True&loc=Local")
	// if err != nil {
	// 	log.Print(err)
	// 	return
	// }

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

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(time.Hour)

	if err := db.Ping(); err != nil {
		log.Print(err)
		return
	}

	var ccs []CronsCrontabs

	err = DB.Find(&ccs).Error
	if err != nil {
		log.Print("find err =", err)
	}

	fmt.Println(ccs)

}
