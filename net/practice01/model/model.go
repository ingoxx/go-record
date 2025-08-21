package model

//公共的model
//定义消息类型常量,根据不同的消息类型分配不同的函数进行处理
const (
	LoginMessAgeType                = "LoginMessAge"
	LoginResultType                 = "LoginResult"
	RegisterType                    = "Register"
	RegisterResultType              = "RegisterResult"
	UserStatusChangeType            = "UserStatusChange"
	ChatMessageType                 = "ChatMessage"
	ChatMessageResultType           = "ChatMessageResult"
	ChatUserToUserMessageType       = "ChatUserToUserMessage"
	ChatUserToUserMessageResultType = "ChatUserToUserMessageResult"
	ChatRrcordMessageResultType     = "ChatRrcordMessageResult"
	FileMessageType                 = "FileMessage"
)

const (
	OnLine  = iota //0
	Offline        //1
	Leisure        //2
)

type UserInfor struct {
	UserId     int    `json:"userid"`
	Passwd     string `json:"passwd"`
	UserStatus int    `json:"userStatus"` //用户状态:在线 离线 空闲

}

//定义消息传输的数据结构,是客户端与服务端传输数据的唯一数据结构
//后面聊天的数据收发都会封装在这个结构体里传输
type Message struct {
	Type string `json:"type"` //消息类型
	Data string `json:"data"` //消息本体如:LoginMessage, RegisterMessage等
}

//登录的结构体,也就是客户端登录提交参数
type LoginMessage struct {
	UserId int    `json:"userId"`
	Passwd string `json:"passwd"`
}

//登录消的响应返回给客户端数据
type LoginResult struct {
	Code     int    `json:"code"`
	UsersId  []int  `json:"usersId "`    //返回在线用户
	AllUsers []int  `json:"offlineList"` //返回所有用户
	Err      string `json:"err"`
}

//注册的结构体,也就是客户端注册提交参数
type RegisterMessage struct {
	UserInfor
}

//注册响应的返回给客户端数据
type RegisterResult struct {
	Code int    `json:"code"`
	Err  string `json:"err"`
}

//服务端推送用户状态变化的消息
type UserStatusChange struct {
	UserId      int   `json:"userid"`
	UserStatus  int   `json:"userStatus"`
	OfflineList []int `json:"offlineList"`
}

//群聊聊消息结构体
type ChatMessage struct {
	UserInfor
	Content string `json:"content"`
}

//群聊消息结构体
type ChatMessageResult struct {
	UserInfor
	Content string `json:"content"`
}

//一对一消息结构体
type ChatUserToUserMessage struct {
	UserInfor
	Content string `json:"content"`
	Sender  int    `json:"sender"`
}

//一对一消息结构体
type ChatUserToUserMessageResult struct {
	UserInfor
	Content string `json:"content"`
}

//聊天记录
type ChatRrcordMessageResult struct {
	MessageList map[int][]ChatUserToUserMessage `json:"MessageList"`
	ChatRecord  []ChatUserToUserMessage
}

//发送文件
type FileMessage struct {
	FileContent []byte
	FileName    string
	UserId      int
	Sender      int
}

type FileMessageResult struct {
	Code int    `json:"code"`
	Err  string `json:"err"`
}
