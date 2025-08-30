package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/importcjj/sensitive"
	"github.com/ingoxx/go-record/http/wx/pkg/config"
	cuerr "github.com/ingoxx/go-record/http/wx/pkg/error"
	"github.com/ingoxx/go-record/http/wx/pkg/eva"
	"github.com/ingoxx/go-record/http/wx/pkg/form"
	"github.com/ingoxx/go-record/http/wx/pkg/redis"
	"github.com/ingoxx/go-record/http/wx/utils/ddw"
	"github.com/ingoxx/go-record/http/wx/utils/openid"
	"github.com/mozillazg/go-pinyin"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var (
	// 字段请求验证器
	validate = validator.New()
	// 脏字库过滤器
	filter = sensitive.New()
)

// Group 一个群聊包含多个客户端连接 + 消息历史
type Group struct {
	Clients  map[*websocket.Conn]bool
	Messages []Message
	Lock     sync.Mutex
}

// Message 用户聊天数据的数据结构
type Message struct {
	GroupID   string `json:"group_id"`
	UserID    string `json:"user_id"`
	NickName  string `json:"nick_name"`
	SenderID  string `json:"sender_id"`
	Content   string `json:"content"`
	Time      string `json:"time"`
	Type      string `json:"type"`
	AvaImg    string `json:"ava_img"`
	UserCount int    `json:"user_count"` // 当前群人数
}

// Resp 响应的数据结构
type Resp struct {
	w         http.ResponseWriter
	OtherData interface{} `json:"other_data"`
	Data      interface{} `json:"data"`
	Msg       string      `json:"msg"`
	Code      int         `json:"code"`
}

func (r Resp) message(rd Resp) ([]byte, error) {
	b, err := json.Marshal(rd)
	if err != nil {
		return b, err
	}
	return b, nil
}

func (r Resp) h(rd Resp) {
	message, err := r.message(rd)
	if err != nil {
		log.Printf("[ERROR] fail to respone, error '%v'", err)
		return
	}
	if _, err := r.w.Write(message); err != nil {
		log.Printf("[ERROR] fail to respone, error '%v'", err)
		return
	}
}

var (
	// 所有群: groupID => Group
	groups   = make(map[string]*Group)
	groupsMu sync.Mutex

	// 全局广播
	broadcast = make(chan Message)

	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
)

func main() {
	log.Println(config.Version)

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", handleConnections)
	mux.HandleFunc("/get-online", handleOnline)
	mux.HandleFunc("/get-join-users", handleGetJoinUsers)
	mux.HandleFunc("/user-join-group", handleUserJoinGroup)
	mux.HandleFunc("/get-all-online-data", handleAllOnlineData)
	mux.HandleFunc("/user-add-square", handleAddSquare)
	mux.HandleFunc("/check-list", handleCheckAddAddrList)
	mux.HandleFunc("/add-square-refuse", handleAddAddrRefuse)
	mux.HandleFunc("/add-square-pass", handleAddAddrPass)
	mux.HandleFunc("/show-square", handleShowSportsSquare)
	mux.HandleFunc("/wx-login", handleWxLogin)
	mux.HandleFunc("/get-all-sports", handleGetAllSports)
	mux.HandleFunc("/wx-upload", handleWxUpload)
	mux.HandleFunc("/get-user-reviews", handleGetUserReviews)
	mux.HandleFunc("/update-sport-reviews", handleUpdateUserReviews)
	mux.HandleFunc("/user-liked-reviews", handleUserLikedReviews)
	mux.HandleFunc("/wx-user-info-update", handleWxUserInfoUpdate)
	mux.HandleFunc("/update-sports-venue", handleUpdateSportsVenue)
	mux.HandleFunc("/get-user-list", handleGetCityList)

	// 启动广播处理器
	go handleBroadcast()

	log.Println("Server started on :11806")
	log.Fatal(http.ListenAndServe(":11806", mux))
}

// 获取所有用户
func handleGetCityList(w http.ResponseWriter, r *http.Request) {
	var rp = Resp{w: w}
	if r.Method != http.MethodGet {
		rp.h(Resp{
			Msg:  "invalid request",
			Code: 1001,
			Data: "0",
		})
		return
	}

	uid := r.FormValue("uid")
	if uid != config.Admin {
		rp.h(Resp{
			Msg:  "invalid parameter",
			Code: 1002,
			Data: "0",
		})
		return
	}

	if err := redis.NewRM().GetWxOpenid(uid); err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1003,
			Data: "0",
		})
		return
	}

	users, err := redis.NewRM().GetAllWxUsers()
	if err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1004,
			Data: "0",
		})
		return
	}

	rp.h(Resp{
		Msg:  "ok",
		Code: 1000,
		Data: users,
	})
}

// handleUpdateSportsVenue 更新场地信息
func handleUpdateSportsVenue(w http.ResponseWriter, r *http.Request) {
	var rp = Resp{w: w}
	if r.Method != http.MethodPost {
		rp.h(Resp{
			Msg:  "invalid request",
			Code: 1001,
			Data: "0",
		})
		return
	}
	uid := r.FormValue("uid")
	if uid == "" {
		rp.h(Resp{
			Msg:  "invalid parameter",
			Code: 1002,
			Data: "0",
		})
		return
	}

	if err := redis.NewRM().GetWxOpenid(uid); err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1003,
			Data: "0",
		})
		return
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1004,
			Data: "0",
		})
		return
	}

	defer r.Body.Close()

	var data *form.UpdateVenueInfo
	if err := json.Unmarshal(b, &data); err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1005,
			Data: "0",
		})
		return
	}

	if err := validate.Struct(data); err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1006,
			Data: "0",
		})
		return
	}

	cityPy := pinyin.LazyPinyin(data.City, pinyin.NewArgs())
	fullKey := fmt.Sprintf("%s_%s", strings.Join(cityPy, ""), data.SportKey) // 拼接的key：shenzhenshi_bks
	data.CityPy = strings.Join(cityPy, "")
	data.SportKey = fullKey

	if !openid.NewWhiteList(uid).IsWhite() {
		if err := ddw.NewDDWarn(fmt.Sprintf("用户id：%s更新了场地图片，场地id：%s，场地类型：%s，场地城市：%s，场地图片：%s", uid, data.Id, data.SportKey, data.City, data.Img)).Send(); err != nil {
			log.Println(err.Error())
		}
	}

	ol, err := redis.NewRM().UpdateVenueInfo(data)
	if err := json.Unmarshal(b, &data); err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1007,
			Data: "0",
		})
		return
	}

	rp.h(Resp{
		Msg:  "ok",
		Code: 1000,
		Data: ol,
	})
}

// handleWxUserInfoUpdate 微信用户信息更新
func handleWxUserInfoUpdate(w http.ResponseWriter, r *http.Request) {
	var rp = Resp{w: w}
	if r.Method != http.MethodPost {
		rp.h(Resp{
			Msg:  "invalid request",
			Code: 1001,
			Data: "0",
		})
		return
	}
	uid := r.FormValue("uid")
	if uid == "" {
		rp.h(Resp{
			Msg:  "invalid parameter",
			Code: 1002,
			Data: "0",
		})
		return
	}

	if err := redis.NewRM().GetWxOpenid(uid); err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1003,
			Data: "0",
		})
		return
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1004,
			Data: "0",
		})
		return
	}

	defer r.Body.Close()

	var data *form.WxOpenidList
	if err := json.Unmarshal(b, &data); err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1005,
			Data: "0",
		})
		return
	}

	data.Time = time.Now().Format("2006-01-02 15:04:05")
	wxOpenid, err := redis.NewRM().SetWxOpenid(data)
	if err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1006,
			Data: "0",
		})
		return
	}

	if !openid.NewWhiteList(data.Openid).IsWhite() {
		if err := ddw.NewDDWarn(fmt.Sprintf("用户id：%s，打开了小程序", data.Openid)).Send(); err != nil {
			log.Println(err.Error())
		}
	}

	rp.h(Resp{
		Msg:       "ok",
		Code:      1000,
		Data:      wxOpenid,
		OtherData: wxOpenid,
	})

}

// handleUserLikedReviews 用户点赞
func handleUserLikedReviews(w http.ResponseWriter, r *http.Request) {
	var rp = Resp{w: w}
	if r.Method != http.MethodPost {
		rp.h(Resp{
			Msg:  "invalid request",
			Code: 1001,
			Data: "0",
		})
		return
	}
	uid := r.FormValue("uid")
	sportKey := r.FormValue("key")
	if uid == "" && sportKey == "" {
		rp.h(Resp{
			Msg:  "invalid parameter",
			Code: 1002,
			Data: "0",
		})
		return
	}

	if err := redis.NewRM().GetWxOpenid(uid); err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1003,
			Data: "0",
		})
		return
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1004,
			Data: "0",
		})
		return
	}
	defer r.Body.Close()

	var data *form.MsgBoard
	if err := json.Unmarshal(b, &data); err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1005,
			Data: "0",
		})
		return
	}

	if err := validate.Struct(data); err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1006,
			Data: "0",
		})
		return
	}

	if !openid.NewWhiteList(data.User).IsWhite() {
		if err := ddw.NewDDWarn(fmt.Sprintf("用户: %s, \n群组：%s, \n点赞了评价：%s\n", data.User, data.GroupId, data.Evaluate)).Send(); err != nil {
			log.Println(err.Error())
		}
	}

	ol, err := redis.NewRM().UserLikedReviews(data, sportKey)
	if err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1007,
			Data: "0",
		})
		return
	}

	rp.h(Resp{
		Msg:  "ok",
		Code: 1000,
		Data: ol,
	})

}

// handleUpdateUserReviews 用户提交对某个场地的评价
func handleUpdateUserReviews(w http.ResponseWriter, r *http.Request) {
	var rp = Resp{w: w}
	if r.Method != http.MethodPost {
		rp.h(Resp{
			Msg:  "invalid request",
			Code: 1001,
			Data: "0",
		})
		return
	}
	uid := r.FormValue("uid")
	sportKey := r.FormValue("key")
	if uid == "" && sportKey == "" {
		rp.h(Resp{
			Msg:  "invalid parameter",
			Code: 1002,
			Data: "0",
		})
		return
	}

	if err := redis.NewRM().GetWxOpenid(uid); err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1003,
			Data: "0",
		})
		return
	}
	b, err := io.ReadAll(r.Body)
	if err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1004,
			Data: "0",
		})
		return
	}
	defer r.Body.Close()

	var data *form.MsgBoard
	if err := json.Unmarshal(b, &data); err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1005,
			Data: "0",
		})
		return
	}

	if err := validate.Struct(data); err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1006,
			Data: "0",
		})
		return
	}

	if !openid.NewWhiteList(data.User).IsWhite() {
		if err := ddw.NewDDWarn(fmt.Sprintf("用户: %s, \n群组：%s, \n提交了评价：%s\n", data.User, data.GroupId, data.Evaluate)).Send(); err != nil {
			log.Println(err.Error())
		}
	}

	if err := filter.LoadWordDict("./dict.txt"); err != nil {
		log.Fatalln("无法读取脏字库文件", err.Error())
	}

	data.Evaluate = filter.Replace(data.Evaluate, '*') // 屏蔽一些不友好的留言
	ol, err := redis.NewRM().UpdateMsgBoard(data, sportKey)
	if err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1007,
			Data: "0",
		})
		return
	}

	rp.h(Resp{
		Msg:  "ok",
		Code: 1000,
		Data: ol,
	})

}

// handleGetUserReviews 获取某个场地的所有用户评价
func handleGetUserReviews(w http.ResponseWriter, r *http.Request) {
	var rp = Resp{w: w}
	if r.Method != http.MethodGet {
		rp.h(Resp{
			Msg:  "invalid request",
			Code: 1001,
			Data: "0",
		})
		return
	}

	gid := r.FormValue("gid")
	uid := r.FormValue("uid")
	sportKey := r.FormValue("key")
	if gid == "" || uid == "" || sportKey == "" {
		rp.h(Resp{
			Msg:  "invalid parameter",
			Code: 1002,
			Data: "0",
		})
		return
	}
	if err := redis.NewRM().GetWxOpenid(uid); err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1003,
			Data: "0",
		})
		return
	}
	ol, err := redis.NewRM().GetMsgBoard(gid, sportKey)
	if err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1004,
			Data: "0",
		})
		return
	}

	rp.h(Resp{
		Msg:  "ok",
		Code: 1000,
		Data: ol,
	})
}

// handleWxUpload 上传文件
func handleWxUpload(w http.ResponseWriter, r *http.Request) {
	var rp = Resp{w: w}
	if r.Method != http.MethodPost {
		rp.h(Resp{
			Msg:  "invalid request",
			Code: 1001,
			Data: "0",
		})
		return
	}

	fileName := r.FormValue("filename")
	uid := r.FormValue("uid")
	up := r.FormValue("user_upload")
	nn := r.FormValue("nick_name")
	if uid == "" || fileName == "" {
		rp.h(Resp{
			Msg:  "invalid parameter",
			Code: 1002,
			Data: "0",
		})
		return
	}

	if err := redis.NewRM().GetWxOpenid(uid); err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1003,
			Data: "0",
		})
		return
	}

	// 限制上传文件大小（例：10MB）
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1003,
			Data: "0",
		})
		return
	}

	// 获取文件
	file, _, err := r.FormFile("file")
	if err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1004,
			Data: "0",
		})
		return
	}

	defer file.Close()

	// 创建保存路径（当前目录的 uploads 文件夹）
	if err := os.MkdirAll(config.UploadPath, os.ModePerm); err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1005,
			Data: "0",
		})
		return
	}

	// 创建目标文件
	dstPath := filepath.Join(config.UploadPath, fileName)
	dst, err := os.Create(dstPath)
	if err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1006,
			Data: "0",
		})
		return
	}
	defer dst.Close()

	// 拷贝内容
	if _, err := io.Copy(dst, file); err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1007,
			Data: "0",
		})
		return
	}

	if !openid.NewWhiteList(uid).IsWhite() {
		if err := ddw.NewDDWarn(fmt.Sprintf("用户: %s, 上传了头像\n", uid)).Send(); err != nil {
			log.Println(err.Error())
		}
	}

	// 更新头像
	if up == "1" {
		var data = &form.WxOpenidList{
			Openid:   uid,
			Img:      fmt.Sprintf("%s/%s", config.ImgUrl, fileName),
			NickName: nn,
			Time:     time.Now().Format("2006-01-02 15:04:05"),
		}

		wxOpenid, err := redis.NewRM().SetWxOpenid(data)
		if err != nil {
			rp.h(Resp{
				Msg:  err.Error(),
				Code: 1008,
				Data: "0",
			})
			return
		}

		rp.h(Resp{
			Msg:       "ok",
			Code:      1000,
			Data:      wxOpenid.Openid,
			OtherData: wxOpenid,
		})

		return
	}

	rp.h(Resp{
		Msg:  "ok",
		Code: 1000,
		Data: "0",
	})

}

// handleUserJoinGroup 用户点击加入某个球局
func handleUserJoinGroup(w http.ResponseWriter, r *http.Request) {
	var rp = Resp{w: w}
	if r.Method != http.MethodPost {
		rp.h(Resp{
			Msg:  "invalid request",
			Code: 1001,
			Data: "0",
		})
		return
	}

	uid := r.FormValue("uid")
	if uid == "" {
		rp.h(Resp{
			Msg:  "invalid parameter",
			Code: 1002,
			Data: "0",
		})
		return
	}

	if err := redis.NewRM().GetWxOpenid(uid); err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1003,
			Data: "0",
		})
		return
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1004,
			Data: "0",
		})
		return
	}

	defer r.Body.Close()

	var data *form.JoinGroupUsers
	if err := json.Unmarshal(b, &data); err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1005,
			Data: "0",
		})
		return
	}

	if err := validate.Struct(data); err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1005,
			Data: "0",
		})
		return
	}

	if !openid.NewWhiteList(data.User).IsWhite() {
		if err := ddw.NewDDWarn(fmt.Sprintf("用户: %s, 群组：%s, 用户点击了组队按钮\n", data.User, data.GroupId)).Send(); err != nil {
			log.Println(err.Error())
		}
	}

	if data.NickName == "" {
		data.NickName = eva.NewSportType("bks").RandomNickname()
	}

	if data.Skill == "" {
		data.Skill = eva.NewSportType("bks").RandomSkill()
	}

	data.Time = time.Now().Format("2006-01-02 15:04:05")

	ol, err := redis.NewRM().JoinGroupUpdate(data)
	if err != nil {
		var dr *cuerr.DuplicateError
		if errors.Is(err, dr) {
			rp.h(Resp{
				Msg:  err.Error(),
				Code: 1006,
				Data: "0",
			})
			return
		}
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1007,
			Data: "0",
		})
		return
	}

	rp.h(Resp{
		Msg:  "ok",
		Code: 1000,
		Data: ol,
	})

}

// handleGetJoinUsers 获取某个组加入/退出的所有用户信息
func handleGetJoinUsers(w http.ResponseWriter, r *http.Request) {
	var rp = Resp{w: w}
	if r.Method != http.MethodGet {
		rp.h(Resp{
			Msg:  "invalid request",
			Code: 1001,
			Data: "0",
		})
		return
	}

	gid := r.FormValue("gid")
	uid := r.FormValue("uid")
	if gid == "" || uid == "" {
		rp.h(Resp{
			Msg:  "invalid parameter",
			Code: 1002,
			Data: "0",
		})
		return
	}

	if err := redis.NewRM().GetWxOpenid(uid); err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1003,
			Data: "0",
		})
		return
	}

	ol, err := redis.NewRM().GetJoinGroupUsers(gid)
	if err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1004,
			Data: "0",
		})
		return
	}

	rp.h(Resp{
		Msg:  "ok",
		Code: 1000,
		Data: ol,
	})
}

// handleGetAllSports 获取所有运动场地类型
func handleGetAllSports(w http.ResponseWriter, r *http.Request) {
	var rp = Resp{w: w}
	if r.Method != http.MethodGet {
		rp.h(Resp{
			Msg:  "invalid request",
			Code: 1001,
			Data: "0",
		})
		return
	}

	uid := r.FormValue("uid")
	if uid == "" {
		rp.h(Resp{
			Msg:  "invalid parameter",
			Code: 1002,
			Data: "0",
		})
		return
	}

	if err := redis.NewRM().GetWxOpenid(uid); err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1003,
			Data: "0",
		})
		return
	}

	sl, err := redis.NewRM().GetSportList()
	if err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1004,
			Data: "0",
		})
		return
	}

	rp.h(Resp{
		Msg:  "ok",
		Code: 1000,
		Data: sl,
	})

}

// handleAllOnlineData 获取所有的在线人数数据
func handleAllOnlineData(w http.ResponseWriter, r *http.Request) {
	var rp = Resp{w: w}
	if r.Method != http.MethodGet {
		rp.h(Resp{
			Msg:  "invalid request",
			Code: 1001,
			Data: "0",
		})
		return
	}

	uid := r.FormValue("uid")
	if uid == "" {
		rp.h(Resp{
			Msg:  "invalid parameter",
			Code: 1002,
			Data: "0",
		})
		return
	}
	if err := redis.NewRM().GetWxOpenid(uid); err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1003,
			Data: "0",
		})
		return
	}
	data, err := redis.NewRM().GetAllOnlineData()
	if err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1004,
			Data: "0",
		})
		return
	}

	rp.h(Resp{
		Msg:  "ok",
		Code: 1000,
		Data: data,
	})

}

// handleWxLogin 微信登陆
func handleWxLogin(w http.ResponseWriter, r *http.Request) {
	var rp = Resp{w: w}
	if r.Method != http.MethodPost {
		rp.h(Resp{
			Msg:  "invalid request",
			Code: 1001,
			Data: "0",
		})
		return
	}
	var codeData map[string]interface{}
	bd, err := io.ReadAll(r.Body)
	if err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1002,
			Data: "0",
		})
		return
	}
	if err = json.Unmarshal(bd, &codeData); err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1003,
			Data: "0",
		})
		return
	}

	v := url.Values{}
	v.Add("appid", config.WxAppid)
	v.Add("secret", config.WxSecret)
	v.Add("js_code", codeData["code"].(string))
	v.Add("grant_type", config.WxAuth)

	urlName := config.WxLoginUrl + v.Encode()
	re, err := http.Get(urlName)
	if err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1004,
			Data: "0",
		})
		return
	}

	defer re.Body.Close()

	b, err := io.ReadAll(re.Body)
	if err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1005,
			Data: "0",
		})
		return
	}

	var data *form.WxOpenidList
	if err := json.Unmarshal(b, &data); err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1006,
			Data: "0",
		})
		return
	}

	data.Time = time.Now().Format("2006-01-02 15:04:05")
	wxOpenid, err := redis.NewRM().SetWxOpenid(data)
	if err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1007,
			Data: "0",
		})
		return
	}

	if !openid.NewWhiteList(data.Openid).IsWhite() {
		if err := ddw.NewDDWarn(fmt.Sprintf("用户id：%s，打开了小程序", data.Openid)).Send(); err != nil {
			log.Println(err.Error())
		}
	}

	rp.h(Resp{
		Msg:       "ok",
		Code:      1000,
		Data:      wxOpenid.Openid,
		OtherData: wxOpenid,
	})

}

// handleShowSportsSquare 根据用户传入的坐标显示用户当前位置附近所有运动场地
func handleShowSportsSquare(w http.ResponseWriter, r *http.Request) {
	var rp = Resp{w: w}
	if r.Method != http.MethodGet {
		rp.h(Resp{
			Msg:  "invalid request",
			Code: 1001,
			Data: "0",
		})
		return
	}

	lng := r.FormValue("lng")
	lat := r.FormValue("lat")
	city := r.FormValue("city") // 中文城市名：深圳市
	uid := r.FormValue("uid")
	sportKey := r.FormValue("sport_key") // 运动类型key：bks等等
	keyWord := r.FormValue("sport_name") // 中文运动场地名称：篮球场，羽毛球场

	if lng == "" || lat == "" || city == "" || uid == "" {
		rp.h(Resp{
			Msg:  "invalid parameter",
			Code: 1002,
			Data: "0",
		})
		return
	}

	// 默认获取的场地是篮球场
	if sportKey == "" || keyWord == "" {
		sportKey = "bks"
		keyWord = "篮球场"
	}

	if err := redis.NewRM().GetWxOpenid(uid); err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1003,
			Data: "0",
		})
		return
	}

	cityPy := pinyin.LazyPinyin(city, pinyin.NewArgs())
	fullKey := fmt.Sprintf("%s_%s", strings.Join(cityPy, ""), sportKey) // 拼接的key：shenzhenshi_bks
	ol, _, err := redis.NewRM().GetAllData(fullKey, city, keyWord, lat, lng, sportKey)
	if err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1004,
			Data: "0",
		})
		return
	}

	if !openid.NewWhiteList(uid).IsWhite() {
		if err := ddw.NewDDWarn(fmt.Sprintf("用户id：%s，城市：%s，选择了：%s运动", uid, city, keyWord)).Send(); err != nil {
			log.Println(err.Error())
		}
	}

	rp.h(Resp{
		Msg:       "ok",
		Code:      1000,
		Data:      "0",
		OtherData: ol,
	})
}

// handleCheckAddAddrList 需要审核的地址列表
func handleCheckAddAddrList(w http.ResponseWriter, r *http.Request) {
	var rp = Resp{w: w}
	if r.Method != http.MethodGet {
		rp.h(Resp{
			Msg:  "invalid request",
			Code: 1001,
			Data: "0",
		})
		return
	}
	uid := r.FormValue("uid")
	if uid != config.Admin {
		rp.h(Resp{
			Msg:  "您有没有权限哟",
			Code: 1002,
			Data: "0",
		})
		return
	}

	if err := redis.NewRM().GetWxOpenid(uid); err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1003,
			Data: "0",
		})
		return
	}

	list, err := redis.NewRM().GetAddrList()
	if err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1004,
			Data: "0",
		})
		return
	}

	rp.h(Resp{
		Msg:  "ok",
		Code: 1000,
		Data: list,
	})
}

// handleAddAddrRefuse 删除不符合要求的用户提交的添加地址请求
func handleAddAddrRefuse(w http.ResponseWriter, r *http.Request) {
	var rp = Resp{w: w}
	if r.Method != http.MethodPost {
		rp.h(Resp{
			Msg:  "invalid request",
			Code: 1001,
			Data: "0",
		})
		return
	}

	uid := r.FormValue("uid")
	if uid != "ogR3E62jXXJMbVcImRqMA1gTSegM" {
		rp.h(Resp{
			Msg:  "invalid parameter",
			Code: 1002,
			Data: "0",
		})
		return
	}
	if err := redis.NewRM().GetWxOpenid(uid); err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1003,
			Data: "0",
		})
		return
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1004,
			Data: "0",
		})
		return
	}

	defer r.Body.Close()

	var data form.PassAddrReqForm

	if err := json.Unmarshal(b, &data); err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1005,
			Data: "0",
		})
		return
	}
	if err := validate.Struct(data); err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1006,
			Data: "0",
		})
		return
	}
	nd, err := redis.NewRM().UpdateAddrList(data.Id, false)
	if err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1007,
			Data: "0",
		})
		return
	}

	rp.h(Resp{
		Msg:  "ok",
		Code: 1000,
		Data: nd,
	})

}

// handleAddAddrPass 通过用户提交的添加地址请求
func handleAddAddrPass(w http.ResponseWriter, r *http.Request) {
	var rp = Resp{w: w}
	if r.Method != http.MethodPost {
		rp.h(Resp{
			Msg:  "invalid request",
			Code: 1001,
			Data: "0",
		})
		return
	}

	uid := r.FormValue("uid")
	if uid != "ogR3E62jXXJMbVcImRqMA1gTSegM" {
		rp.h(Resp{
			Msg:  "invalid parameter",
			Code: 1002,
			Data: "0",
		})
		return
	}

	if err := redis.NewRM().GetWxOpenid(uid); err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1003,
			Data: "0",
		})
		return
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1004,
			Data: "0",
		})
		return
	}

	defer r.Body.Close()

	var data form.PassAddrReqForm

	if err := json.Unmarshal(b, &data); err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1005,
			Data: "0",
		})
		return
	}

	if err := validate.Struct(data); err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1006,
			Data: "0",
		})
		return
	}

	// 更新已经存在的场地信息
	if data.UpdateType == "2" {
		ud := &form.UpdateVenueInfo{
			SportKey: data.City,
			Id:       data.Id,
			Img:      data.Img,
		}
		if _, err := redis.NewRM().UpdateVenueInfo(ud); err != nil {
			rp.h(Resp{
				Msg:  err.Error(),
				Code: 1007,
				Data: "0",
			})
			return
		}

	}

	// key:shenzhenshi_bks
	if _, err := redis.NewRM().Update(data.City, data.Id, data.UpdateType); err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1008,
			Data: "0",
		})
		return
	}

	nd, err := redis.NewRM().UpdateAddrList(data.Id, true)
	if err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1009,
			Data: "0",
		})
		return
	}

	rp.h(Resp{
		Msg:  "ok",
		Code: 1000,
		Data: nd,
	})

}

// handleAddSquare 用户提交的地址添加到地址列表
func handleAddSquare(w http.ResponseWriter, r *http.Request) {
	var rp = Resp{w: w}
	if r.Method != http.MethodPost {
		rp.h(Resp{
			Msg:  "invalid request",
			Code: 1001,
			Data: "0",
		})
		return
	}

	uid := r.FormValue("uid")
	if uid == "" {
		rp.h(Resp{
			Msg:  "invalid parameter",
			Code: 1002,
			Data: "0",
		})
		return
	}

	if err := redis.NewRM().GetWxOpenid(uid); err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1003,
			Data: "0",
		})
		return
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1004,
			Data: "0",
		})
		return
	}

	defer r.Body.Close()

	var data form.AddrListForm

	if err := json.Unmarshal(b, &data); err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1005,
			Data: "0",
		})
		return
	}

	if err := validate.Struct(data); err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1006,
			Data: "0",
		})
		return
	}

	data.CityPy = strings.Join(pinyin.LazyPinyin(data.City, pinyin.NewArgs()), "")
	data.SportKey = fmt.Sprintf("%s_%s", data.CityPy, data.SportKey) // 拼接的key：shenzhenshi_bks
	if err := redis.NewRM().UserAddAddrReq(data); err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1007,
			Data: "0",
		})
		return
	}

	if !openid.NewWhiteList(uid).IsWhite() {
		if err := ddw.NewDDWarn(fmt.Sprintf("用户id：%s，城市：%s, 添加地址：%s", uid, data.City, data.Addr)).Send(); err != nil {
			log.Println(err.Error())
		}
	}

	rp.h(Resp{
		Msg:  "ok",
		Code: 1000,
		Data: "0",
	})

}

// handleOnline 群组的在线人数
func handleOnline(w http.ResponseWriter, r *http.Request) {
	var rp = Resp{w: w}
	if r.Method != http.MethodGet {
		rp.h(Resp{
			Msg:  "invalid request",
			Code: 1001,
			Data: "0",
		})
		return
	}

	gid := r.FormValue("gid")
	uid := r.FormValue("uid")
	if gid == "" || uid == "" {
		rp.h(Resp{
			Msg:  "invalid parameter",
			Code: 1002,
			Data: "0",
		})
		return
	}

	if err := redis.NewRM().GetWxOpenid(uid); err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1003,
			Data: "0",
		})
		return
	}

	ol, err := redis.NewRM().GetGroupOnline(gid)
	if err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1004,
			Data: "0",
		})
		return
	}

	rp.h(Resp{
		Msg:  "ok",
		Code: 1000,
		Data: ol,
	})

}

// handleConnections 聊天室接收到的用户发的信息
func handleConnections(w http.ResponseWriter, r *http.Request) {
	uid := r.FormValue("uid")
	if uid == "" {
		http.Error(w, "Forbidden: Invalid UID", http.StatusForbidden)
		return
	}

	if err := redis.NewRM().GetWxOpenid(uid); err != nil {
		http.Error(w, "Forbidden: Invalid UID", http.StatusForbidden)
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket 升级失败:", err)
		return
	}

	if err := filter.LoadWordDict("./dict.txt"); err != nil {
		log.Fatalln("无法读取脏字库文件", err.Error())
	}

	defer ws.Close()

	// 先读取第一条消息，拿到 groupID
	_, msgBytes, err := ws.ReadMessage()
	if err != nil {
		log.Println("初始化连接读取失败:", err)
		return
	}

	var initMsg Message
	if err := json.Unmarshal(msgBytes, &initMsg); err != nil {
		log.Println("解析初始化消息失败:", err)
		return
	}

	// 用户进入了聊天室就发送通知
	if !openid.NewWhiteList(initMsg.UserID).IsWhite() {
		if err := ddw.NewDDWarn(fmt.Sprintf("用户: %s, 进入了群组：%s聊天室\n", initMsg.UserID, initMsg.GroupID)).Send(); err != nil {
			log.Println(err.Error())
		}
	}

	groupID := initMsg.GroupID

	// 把连接放入对应群
	groupsMu.Lock()
	if _, ok := groups[groupID]; !ok {
		groups[groupID] = &Group{
			Clients:  make(map[*websocket.Conn]bool),
			Messages: []Message{},
		}
	}
	group := groups[groupID]
	groupsMu.Unlock()

	// 给群组添加ws
	group.Lock.Lock()
	group.Clients[ws] = true
	userCount := len(group.Clients)
	if groupID != "" {
		gn := fmt.Sprintf("%s_%s", config.OnlineKey, groupID)
		if err := redis.NewRM().Set(gn, userCount, time.Second*time.Duration(7200)); err != nil {
			log.Printf("[ERROR] 写入redis失败, 错误信息：%v", err)
		}
	}
	group.Lock.Unlock()

	log.Printf("用户 %s 加入群 %s，当前人数: %d", initMsg.UserID, groupID, userCount)

	// 广播新的群人数
	broadcast <- Message{
		GroupID:   groupID,
		Type:      "count",
		UserCount: userCount,
	}

	// 先把历史消息发给新连接（可选）
	group.Lock.Lock()
	for _, oldMsg := range group.Messages {
		oldMsg.Content = filter.Replace(oldMsg.Content, '*')
		if err := ws.WriteJSON(oldMsg); err != nil {
			log.Println("发送历史消息失败:", err)
		}
	}
	group.Lock.Unlock()

	// 持续读取新消息
	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			// 从群里删除这个连接
			group.Lock.Lock()
			delete(group.Clients, ws)
			userCount = len(group.Clients)
			log.Printf("用户：%s, 组：%s, 当前人数: %d,  断开连接", initMsg.UserID, initMsg.GroupID, userCount)
			if msg.GroupID != "" {
				gn := fmt.Sprintf("%s_%s", config.OnlineKey, msg.GroupID)
				if err := redis.NewRM().Set(gn, userCount, time.Second*time.Duration(7200)); err != nil {
					log.Printf("[ERROR] 写入redis失败, 错误信息：%v", err)
				}
			}
			group.Lock.Unlock()

			// 广播新的群人数
			broadcast <- Message{
				GroupID:   groupID,
				Type:      "count",
				UserCount: userCount,
			}

			break
		}

		log.Printf("用户: %s, 群组：%s, 发送的内容: %s\n", msg.UserID, msg.GroupID, msg.Content)
		if !openid.NewWhiteList(msg.UserID).IsWhite() {
			if err := ddw.NewDDWarn(fmt.Sprintf("用户: %s, 群组：%s, 发送的内容: %s\n", msg.UserID, msg.GroupID, msg.Content)).Send(); err != nil {
				log.Println(err.Error())
			}
		}

		// 普通消息
		msg.Type = "normal"

		// 保存到历史
		group.Lock.Lock()
		group.Messages = append(group.Messages, msg)
		group.Lock.Unlock()

		// 广播消息
		broadcast <- msg
	}
}

// handleBroadcast 将用户发送的信息过滤处理后返回到前端
func handleBroadcast() {
	for {
		msg := <-broadcast
		msg.Content = filter.Replace(msg.Content, '*')
		groupID := msg.GroupID

		groupsMu.Lock()
		group, ok := groups[groupID]
		if msg.GroupID != "" {
			gn := fmt.Sprintf("%s_%s", config.OnlineKey, msg.GroupID)
			if err := redis.NewRM().Set(gn, msg.UserCount, time.Second*time.Duration(7200)); err != nil {
				log.Println("[ERROR] fail to save user count.")
			}
		}
		groupsMu.Unlock()
		if !ok {
			continue
		}

		group.Lock.Lock()
		for client := range group.Clients {
			if err := client.WriteJSON(msg); err != nil {
				log.Println("广播失败，删除连接:", err)
				client.Close()
				delete(group.Clients, client)
			}
		}
		group.Lock.Unlock()
	}
}
