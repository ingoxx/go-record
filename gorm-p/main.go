package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/ingoxx/go-record/gorm-p/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
	"time"
	// _ "github.com/go-sql-driver/mysql"
	// "github.com/jinzhu/gorm"
)

var (
	DB     *gorm.DB
	rdPool *redis.Client
)

type OperateLogModel struct {
	gorm.Model
	Url      string    `json:"url" gorm:"type:text;not null"`
	Operator string    `json:"operator" gorm:"not null"`
	Ip       string    `json:"ip" gorm:"not null"`
	Start    time.Time `json:"start" gorm:"-"`
	End      time.Time `json:"end" gorm:"-"`
}

type User struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	Name      string         `json:"name" gorm:"unique;not null"`
	Email     string         `json:"email" gorm:"unique;not null"`
	Hobby     string         `json:"-" gorm:"default:'basketball'"`
	Tel       int            `json:"tel" gorm:"default:168888"`
	Password  string         `json:"password,omitempty" gorm:"not null"`
	Roles     []Role         `json:"roles" gorm:"many2many:role_users"`
	Isopenga  uint           `json:"isopenga" gorm:"default:1"`
	Isopenqr  uint           `json:"isopenqr" gorm:"default:1"`
}

type Permission struct {
	gorm.Model
	Path     string `json:"path" gorm:"not null"`
	Title    string `json:"title" gorm:"not null"`
	Hidden   uint   `json:"hidden" gorm:"default:1;comment:'1是可见,2是隐藏'"`
	ParentId uint   `json:"parentId" gorm:"default 0"`
	Level    uint   `json:"level" gorm:"not null"`
	Roles    []Role `json:"roles" gorm:"many2many:role_permissions"`
}

type Role struct {
	gorm.Model
	RoleName    string `json:"rolename" gorm:"unique;not null"`
	Description string `json:"description" gorm:"-"`
	//不同步更新permission表
	//Permission  []Permission `json:"permission" gorm:"many2many:role_permissions;association_autoupdate:false;association_autocreate:false"`
	//同步更新permission表
	Permissions []Permission `json:"permissions" gorm:"many2many:role_permissions"`
	Users       []User       `json:"users" gorm:"many2many:role_users"`
}

// AssetsModel 服务器列表
type AssetsModel struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Ip          string    `json:"ip" gorm:"not null;unique"`
	Project     string    `json:"project" gorm:"not null"`
	Status      uint      `json:"status" gorm:"default:200;comment:100-服务器异常,200-服务器正常"`
	Operator    string    `json:"operator" gorm:"default:lxb"`
	RamUsage    uint      `json:"ram_usage" gorm:"default:1"`
	DiskUsage   uint      `json:"disk_usage" gorm:"default:1"`
	CpuUsage    uint      `json:"cpu_usage"  gorm:"default:1"`
	Start       time.Time `json:"start" gorm:"default:CURRENT_TIMESTAMP;nullable"`
	User        string    `json:"user" gorm:"default:root"`
	Password    string    `json:"-" gorm:"not null"`
	Key         string    `json:"-" gorm:"type:text"`
	Port        uint      `json:"port" gorm:"default:22"`
	ConnectType uint      `json:"connect_type" gorm:"default:1;comment:1-密码登陆, 2-秘钥登陆"`
}

func main() {
	if err := InitPoolMysql(); err != nil {
		log.Fatalln(err)
	}
	initPoolRedis()
	del_users()
}

func getUser() {
	var u User
	DB.Model(&User{}).Where("id = ?", 9).Find(&u)
	fmt.Println(u)
}

func getServerKey() {
	var am AssetsModel
	DB.Model(&AssetsModel{}).Where("id = ?", 11).Find(&am)
	fmt.Println(am)
}

func countLog() {
	var dataList = make([]map[string]interface{}, 0)
	rows, _ := DB.Raw(`
		SELECT DATE(created_at) as date, count(1) as login_num FROM operate_log_models 
		where DATE(created_at) > NOW() - INTERVAL 7 DAY and url like '%/login%'
		GROUP BY DATE(created_at);
	`).Rows()
	// 遍历结果并填充到 map 中
	for rows.Next() {
		var date string
		var loginNum int
		var data = make(map[string]interface{})

		if err := rows.Scan(&date, &loginNum); err != nil {
			log.Fatal(err)
		}

		parsedTime, _ := time.Parse(time.RFC3339, date)

		data[parsedTime.Format("2006-01-02")] = loginNum
		dataList = append(dataList, data)
	}

	// 输出结果
	for k, v := range dataList {
		fmt.Println(k, v)
	}

}

func find_user() {
	var user = User{Name: "admin"}
	var users []User
	if err := DB.Where(&user).Preload("Roles").Find(&users).Error; err != nil {
		log.Fatalln(err)
	}

	fmt.Println(users)
}

func del_perms() {
	pid := []uint{38}
	var perm Permission
	tx := DB.Begin()
	for _, v := range pid {
		if err := DB.Where("id = ?", v).Find(&perm).Error; err != nil {
			log.Fatalln("err111 >>> ", err)
		}

		if err := tx.Model(&perm).Association("Roles").Clear(); err != nil {
			tx.Rollback()
			log.Fatalln("err222 >>> ", err)
		}

		if err := tx.Unscoped().Delete(&perm, v).Error; err != nil {
			tx.Rollback()
			log.Fatalln("err333 >>> ", err)
		}

	}

	if err := tx.Commit().Error; err != nil {
		log.Fatalln(err)
	}

}

func get_perms() {
	var perms []Permission
	DB.Where("id = 37").Find(&perms)
	fmt.Println(perms)
}

func give_perms_to_role() (err error) {
	var role Role
	var perms []Permission

	if err = DB.Where("id > 0").Find(&perms).Error; err != nil {
		return
	}

	if len(perms) == 0 {
		return errors.New("权限列表空")
	}

	if err = DB.Where("id = ?", 1).Find(&role).Error; err != nil {
		return
	}

	if err = DB.Model(&role).Association("Permissions").Replace(&perms); err != nil {
		return
	}

	//if err = DB.Model(&role).Association("Permissions").Find(&perms); err != nil {
	//	return
	//}
	//
	//fmt.Println(perms)

	return
}

func clear_all_perms() (err error) {
	var role Role
	if err = DB.Where("id = ?", 1).Find(&role).Error; err != nil {
		return
	}

	if err = DB.Model(&role).Association("Permissions").Clear(); err != nil {
		return
	}

	return
}

func check_user_perms(uid uint) {
	var user User
	if err := DB.Preload("Roles.Permissions").First(&user, uid).Error; err != nil {
		log.Fatalln(err)
	}

	permissionsMap := make(map[uint]Permission) // 使用 map 去重
	for _, role := range user.Roles {
		for _, permission := range role.Permissions {
			permissionsMap[permission.ID] = permission
		}
	}

	permissions := make([]Permission, 0, len(permissionsMap))
	for _, permission := range permissionsMap {
		permissions = append(permissions, permission)
	}

	fmt.Println("用户的所有权限：", permissions)
}

func check_user_roles() {
	var user User
	DB.Preload("Roles").Find(&user)
	fmt.Println(user)
}

func get_all_roles() {
	var roles []Role
	DB.Find(&roles)
	fmt.Println(roles)
}

func check_role_perms(rid uint) {
	var role Role
	var perms []Permission
	//DB.Preload("Permissions").Find(&role, rid)
	DB.First(&role, "id = ?", rid)
	DB.Model(&role).Where("id IN ?", []uint{1}).Association("Permissions").Delete(&perms)
	fmt.Println(perms)
}

func getAllPerms() []Permission {
	var perms []Permission
	DB.Where("id > 0").Find(&perms)
	return perms
}

func create_user() {
	rows, _ := DB.Model(&Permission{}).Where("id > ?", 0).Rows()
	defer rows.Close()

	var pdl []Permission

	for rows.Next() {
		var perm Permission
		DB.ScanRows(rows, &perm)
		pdl = append(pdl, perm)
	}

	var rl []Role

	DB.Model(&Role{}).Where("id = ?", 1).Find(&rl)

	fmt.Println(rl)

	//r := Role{
	//	RoleName:    "管理员",
	//	Description: "拥有所有权限",
	//	Permissions: pdl,
	//}
	//
	//rl = append(rl, r)

	u := User{
		Name:     "admin",
		Roles:    rl,
		Password: "123321",
		Email:    "1354198737@qq.com",
	}

	DB.Create(&u)
	DB.Save(&u)
}

func DeleteUser(uid []uint) (err error) {
	var us []User
	tx := DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err = tx.Where("id IN ?", uid).Find(&us).Error; err != nil {
		tx.Rollback()
		return
	}

	if err = tx.Select(clause.Associations).Unscoped().Delete(&us).Error; err != nil {
		tx.Rollback()
		return
	}

	return tx.Commit().Error
}

func getUserName() {
	var u = new(User)
	if err := DB.Where("name = ?", "admin").Preload("Roles").Find(u).Error; err != nil {
		return
	}

	marshal, _ := json.Marshal(u)

	rdPool.Set("admin-rc", marshal, time.Second*2592000)

	fmt.Println(u)
}

func del_users() {
	var us []User
	tx := DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Where("id IN ?", []uint{19}).Find(&us).Error; err != nil {
		tx.Rollback()
		return
	}

	if err := tx.Select(clause.Associations).Unscoped().Delete(&us).Error; err != nil {
		tx.Rollback()
		return
	}

	fmt.Println("OK")
}

func getAllUsers() {
	var us []User
	DB.Debug().Find(&us)
	fmt.Println(us)
}

func InitPoolMysql() (err error) {
	DB, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       config.MyConAddre, // DSN data source name
		DefaultStringSize:         256,               // string 类型字段的默认长度
		DisableDatetimePrecision:  true,              // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,              // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,              // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,             // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{})
	if err != nil {
		return
	}

	sqlDB, err := DB.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	if err != nil {
		return
	}

	return sqlDB.Ping()
}

func initPoolRedis() { //初始化
	rdPool = redis.NewClient(&redis.Options{
		Addr:         config.RedisConAddre,
		DB:           config.RedisUserDb,
		MinIdleConns: 5,
		Password:     config.RedisPwd,
		PoolSize:     5,
		PoolTimeout:  30 * time.Second,
		DialTimeout:  1 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	})
}

func DeleteRole(rid []uint) (err error) {
	tx := DB.Begin()

	return tx.Commit().Error
}
