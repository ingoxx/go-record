package model

import "errors"

var (
	ErrorUserNotExists = errors.New("用户不存在")
	ErrorUserExists    = errors.New("用户已存在,请换个名")
	ErrorUserPwd       = errors.New("密码错误")
	ErrorRegister      = errors.New("注册失败,redis操作错误")
	ErrorJsonMarshal   = errors.New("序列化失败")
	ErrorJsonUnmarshal = errors.New("反序列化失败")
	ErrorFileNotExists = errors.New("文件不存在")
)
