package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/importcjj/sensitive"
	"github.com/ingoxx/go-record/http/wx/form"
	"github.com/ingoxx/go-record/http/wx/redis"
	"github.com/mozillazg/go-pinyin"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	validate = validator.New()
	filter   = sensitive.New()
)

// Group 一个群聊包含多个客户端连接 + 消息历史
type Group struct {
	Clients  map[*websocket.Conn]bool
	Messages []Message
	Lock     sync.Mutex
}

type Message struct {
	GroupID   string `json:"group_id"`
	UserID    string `json:"user_id"`
	Content   string `json:"content"`
	Time      string `json:"time"`
	Type      string `json:"type"`       // normal / count
	UserCount int    `json:"user_count"` // 当前群人数
}

type Resp struct {
	w    http.ResponseWriter
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
	Code int         `json:"code"`
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
	log.Println("version: v1.1.19")

	http.HandleFunc("/ws", handleConnections)
	http.HandleFunc("/get-online", handleOnline)
	http.HandleFunc("/user-add-square", handleAddBasketballSquare)
	http.HandleFunc("/check-list", handleCheckAddAddrList)
	http.HandleFunc("/add-square-refuse", handleAddAddrRefuse)
	http.HandleFunc("/add-square-pass", handleAddAddrPass)
	http.HandleFunc("/show-square", handleShowBasketballSquare)
	http.HandleFunc("/wx-login", handleWxLogin)

	// 启动广播处理器
	go handleBroadcast()

	log.Println("Server started on :11806")
	log.Fatal(http.ListenAndServe(":11806", nil))
}

func handleWxLogin(w http.ResponseWriter, r *http.Request) {
	var rp = Resp{w: w}
	if r.Method != "POST" {
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
			Msg:  "invalid request",
			Code: 1003,
			Data: "0",
		})
		return
	}
	if err = json.Unmarshal(bd, &codeData); err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1004,
			Data: "0",
		})
		return
	}

	v := url.Values{}
	v.Add("appid", "wxbb1377eff3149db4")
	v.Add("secret", "6bae191e5e03aa4cf5731478ab513624")
	v.Add("js_code", codeData["code"].(string))
	v.Add("grant_type", "authorization_code")

	urlName := "https://api.weixin.qq.com/sns/jscode2session?" + v.Encode()
	re, err := http.Get(urlName)
	if err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1005,
			Data: "0",
		})
		return
	}

	defer re.Body.Close()

	b, err := io.ReadAll(re.Body)
	if err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1006,
			Data: "0",
		})
		return
	}

	var wxOpenid struct {
		Openid string `json:"openid"`
	}

	if err := json.Unmarshal(b, &wxOpenid); err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1006,
			Data: "0",
		})
		return
	}

	rp.h(Resp{
		Msg:  "ok",
		Code: 1000,
		Data: wxOpenid.Openid,
	})

}

// handleShowBasketballSquare 根据用户传入的坐标显示用户当前位置附近所有篮球场
func handleShowBasketballSquare(w http.ResponseWriter, r *http.Request) {
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
	city := r.FormValue("city") // 中文

	if lng == "" || lat == "" || city == "" {
		rp.h(Resp{
			Msg:  "invalid parameter",
			Code: 1002,
			Data: "0",
		})
		return
	}
	cityPy := pinyin.LazyPinyin(city, pinyin.NewArgs())
	log.Printf("中文： %s, 拼音： %s\n", city, cityPy)
	ol, err := redis.NewRM().GetAllData(strings.Join(cityPy, ""), city)
	if err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1003,
			Data: ol,
		})
		return
	}

	rp.h(Resp{
		Msg:  "ok",
		Code: 1000,
		Data: ol,
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
	if uid != "ogR3E62jXXJMbVcImRqMA1gTSegM" {
		rp.h(Resp{
			Msg:  "您有没有权限哟",
			Code: 1002,
			Data: "0",
		})
		return
	}

	list, err := redis.NewRM().GetAddrList()
	if err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1003,
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

	b, err := io.ReadAll(r.Body)
	if err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1002,
			Data: "0",
		})
		return
	}

	defer r.Body.Close()

	var data form.PassAddrReqForm

	if err := json.Unmarshal(b, &data); err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1003,
			Data: "0",
		})
		return
	}
	if err := validate.Struct(data); err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1004,
			Data: "0",
		})
		return
	}
	nd, err := redis.NewRM().UpdateAddrList(data.Id)
	if err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1005,
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

// handleAddAddrPass 审核通过用户提交的添加地址请求
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

	b, err := io.ReadAll(r.Body)
	if err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1002,
			Data: "0",
		})
		return
	}

	defer r.Body.Close()

	var data form.PassAddrReqForm

	if err := json.Unmarshal(b, &data); err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1003,
			Data: "0",
		})
		return
	}

	if err := validate.Struct(data); err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1004,
			Data: "0",
		})
		return
	}

	if _, err := redis.NewRM().Update(data.City, data.Id); err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1005,
			Data: "0",
		})
		return
	}

	nd, err := redis.NewRM().UpdateAddrList(data.Id)
	if err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1006,
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

// handleAddBasketballSquare 将用户填入的地址添加到列表中
func handleAddBasketballSquare(w http.ResponseWriter, r *http.Request) {
	var rp = Resp{w: w}
	if r.Method != http.MethodPost {
		rp.h(Resp{
			Msg:  "invalid request",
			Code: 1001,
			Data: "0",
		})
		return
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1002,
			Data: "0",
		})
		return
	}

	defer r.Body.Close()

	var data form.AddAddrForm

	if err := json.Unmarshal(b, &data); err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1003,
			Data: "0",
		})
		return
	}

	if err := validate.Struct(data); err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1004,
			Data: "0",
		})
		return
	}

	data.CityPy = strings.Join(pinyin.LazyPinyin(data.City, pinyin.NewArgs()), "")
	if err := redis.NewRM().UserAddAddrReq(data); err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1005,
			Data: "0",
		})
		return
	}

	rp.h(Resp{
		Msg:  "ok",
		Code: 1000,
		Data: "0",
	})

}

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
	if gid == "" {
		rp.h(Resp{
			Msg:  "invalid parameter",
			Code: 1002,
			Data: "0",
		})
		return
	}

	ol, err := redis.NewRM().Get(gid)
	if err != nil {
		rp.h(Resp{
			Msg:  fmt.Sprintf("'%s' not found", gid),
			Code: 1003,
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

func handleConnections(w http.ResponseWriter, r *http.Request) {
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
	if groupID == "" {
		if err := redis.NewRM().Set(groupID, userCount, 0); err != nil {
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
				if err := redis.NewRM().Set(msg.GroupID, userCount, 0); err != nil {
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

		log.Printf("用户: %s, 发送的内容: %s\n", msg.UserID, msg.Content)

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

func handleBroadcast() {
	for {
		msg := <-broadcast
		msg.Content = filter.Replace(msg.Content, '*')
		groupID := msg.GroupID

		groupsMu.Lock()
		group, ok := groups[groupID]
		if msg.GroupID != "" {
			if err := redis.NewRM().Set(msg.GroupID, msg.UserCount, 0); err != nil {
				log.Println("[ERROR] fail to save user count.")
			}
		}
		groupsMu.Unlock()
		if !ok {
			continue
		}

		group.Lock.Lock()
		for client := range group.Clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Println("广播失败，删除连接:", err)
				client.Close()
				delete(group.Clients, client)
			}
		}
		group.Lock.Unlock()
	}
}
