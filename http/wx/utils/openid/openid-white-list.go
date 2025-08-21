package openid

import "github.com/ingoxx/go-record/http/wx/pkg/config"

type WhiteList struct {
	uid string
}

func NewWhiteList(uid string) WhiteList {
	return WhiteList{
		uid: uid,
	}
}

func (w WhiteList) IsWhite() bool {
	for _, u := range config.Wl {
		if u == w.uid {
			return true
		}
	}

	return false
}
