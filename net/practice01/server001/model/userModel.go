package model

//用户信息字段,用户的增删改查都需要用到
type UserInfo struct {
	UserID int    `json:"UserID"`
	Passwd string `json:"Passwd"`
}
