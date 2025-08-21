package model

import (
	"net"

	"github.com/ingoxx/Golang-practise/net/practice01/model"
)

type CurrentUser struct {
	model.UserInfor
	Conn net.Conn //客户端自己也要维护一个Conn
}
