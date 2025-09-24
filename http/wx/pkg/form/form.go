package form

// BdAPIResponse 百度地图
type BdAPIResponse struct {
	Data    []BdResponse `json:"results"` // 这是一个 DataItem 的切片(数组)
	Message string       `json:"message"`
	Count   int          `json:"total"`
	Status  int          `json:"status"`
}

// BdResponse 百度地图
type BdResponse struct {
	ID       string   `json:"uid"`
	Address  string   `json:"address"`
	AdCode   string   `json:"adcode"`
	PName    string   `json:"pname"`
	City     string   `json:"city"`
	Area     string   `json:"area"`
	Town     string   `json:"town"`
	Name     string   `json:"name"` // 场地简称如: 新湖篮球馆
	Location Location `json:"location"`
}

// GdAPIResponse 高德
type GdAPIResponse struct {
	Data     []GdResponse `json:"pois"` // 这是一个 DataItem 的切片(数组)
	Message  string       `json:"message"`
	Info     string       `json:"info"`
	Count    string       `json:"count"`
	Status   string       `json:"status"`
	InfoCode string       `json:"infocode"`
}

// GdResponse 高德
type GdResponse struct {
	Type         string      `json:"type"`
	Photos       []Photo     `json:"photos"`
	AdName       string      `json:"adname"`
	Tel          interface{} `json:"tel"`
	ID           string      `json:"id"`
	Timestamp    string      `json:"timestamp"`
	Address      interface{} `json:"address"`
	AdCode       string      `json:"adcode"`
	PName        string      `json:"pname"`
	CityName     string      `json:"cityname"`
	BusinessArea interface{} `json:"business_area"`
	Name         string      `json:"name"`
	Location     interface{} `json:"location"`
}

// Photo 高德
type Photo struct {
	//Title []interface{} `json:"title"`
	URL string `json:"url"`
}

// TxAPIResponse 腾讯
type TxAPIResponse struct {
	Status  int          `json:"status"`
	Message string       `json:"message"`
	Data    []TxDataItem `json:"data"` // 这是一个 TxDataItem 的切片(数组)
	Count   int          `json:"count"`
}

// TxDataItem 腾讯
type TxDataItem struct {
	Location Location `json:"location"` // 嵌套了 Location 结构体
	ID       string   `json:"id"`
	Title    string   `json:"title"`
	Address  string   `json:"address"`
	Province string   `json:"province"`
	City     string   `json:"city"`
	District string   `json:"district"`
	// 其他字段如果不需要可以不定义
}

// Location 经纬度
type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

// SaveInRedis 统一的写入格式
type SaveInRedis struct {
	JoinUsers        []*JoinGroupUsers `json:"join_users"`   // 某个运动场地，用户点击加入组件的人数
	UserReviews      []*MsgBoard       `json:"user_reviews"` // 某个场地的所有评价
	VenueUpdateUsers []*AddrListForm   `json:"venue_update_users"`
	//OnlineData             []*OnlineData     `json:"online_data"`
	Tags                   []string `json:"tags"`
	Images                 []string `json:"images"`
	Id                     string   `json:"id"`
	Img                    string   `json:"img"`
	Addr                   string   `json:"addr"`
	Title                  string   `json:"title"`
	UserId                 string   `json:"user_id"`
	Online                 string   `json:"online"`
	Distance               string   `json:"distance"`
	Aid                    string   `json:"aid"` // 接口返回的地址唯一id，再次请求接口返回的id是一致的，更新的时候有用
	JoinUserCount          int      `json:"join_user_count"`
	UserReviewsCount       int      `json:"user_reviews_count"`
	VenueUpdateUsersCount  int      `json:"venue_update_users_count"`
	Lng                    float64  `json:"lng"`
	Lat                    float64  `json:"lat"`
	DisVal                 float64  `json:"dis_val"`
	IsShow                 bool     `json:"is_show"`
	IsShowUserReviews      bool     `json:"is_show_user_reviews"`
	IsShowJoinUsers        bool     `json:"is_show_join_users"`
	IsShowVenueUpdateUsers bool     `json:"is_show_venue_update_users"`
}

type AddrListForm struct {
	UserImg    string  `json:"user_img"` // 用户头像
	Content    string  `json:"content"`  // 更新内容,目前只能统一更新图片,这里都写: 更新了场地图片
	NickName   string  `json:"nick_name"`
	Tags       string  `json:"tags"  validate:"required"`
	Id         string  `json:"id" validate:"required"`
	Addr       string  `json:"addr" validate:"required"`
	UserId     string  `json:"user_id" validate:"required"`      // 添加场地的用户id
	City       string  `json:"city"  validate:"required"`        // 前端传入的是中文
	CityPy     string  `json:"city_py"`                          // 前端传入的中文转成拼音
	SportKey   string  `json:"sport_key" validate:"required"`    // 运动分类，篮球：shenzhenshi_bks,足球：shenzhenshi_fbs...
	UpdateType string  `json:"update_type"  validate:"required"` // 更新类型：1.用户添加的新场地，2.用户更新了场地
	Aid        string  `json:"aid"`                              // api返回的场地的唯一id，就是再次请求返回的id都是一样的
	Img        string  `json:"img"`                              // 场地图片
	Time       string  `json:"time"`                             // 更新时间
	Lat        float64 `json:"lat"`
	Lng        float64 `json:"lng"`
	IsRecord   bool    `json:"is_record"` // true：已记录（审核通过），false：未记录（还未审核通过）
	IsShow     bool    `json:"is_show"`   // 审核列表中的数据，true：隐藏，false：不隐藏
}

type PassAddrReqForm struct {
	Id         string `json:"id" validate:"required"`           // 场地的id
	City       string `json:"city"  validate:"required"`        // 城市，中文
	UpdateType string `json:"update_type"  validate:"required"` // 更新类型：1.用户添加的新场地，2.用户更新了场地
	Img        string `json:"img" validate:"required"`          // 场地图片
}

type WxOpenidList struct {
	Openid   string `json:"openid"`
	Img      string `json:"img"`
	NickName string `json:"nick_name"`
	Time     string `json:"time"`
	City     string `json:"city"`
}

type Sports struct {
	Name    string `json:"name"`
	Key     string `json:"key"`
	Id      string `json:"id"`
	Img     string `json:"img"`
	Checked bool   `json:"checked"`
}

type SportList struct {
	Name     string `json:"name"`
	Key      string `json:"key"`
	Icon     string `json:"icon"`
	Img      string `json:"img"`
	Title    string `json:"title"`
	SportImg string `json:"sport_img"`
	Checked  bool   `json:"checked"`
}

type GroupOnlineStatus struct {
	GroupId    string `json:"group_id"`
	VenueName  string `json:"venue_name"`
	OnlineUser string `json:"online_user"`
}

// JoinGroupUsers 单个用户加入记录
type JoinGroupUsers struct {
	GroupId   string `json:"group_id" validate:"required"`
	User      string `json:"user"`
	NickName  string `json:"nick_name"`
	Img       string `json:"img"`
	Style     string `json:"style"`
	Skill     string `json:"skill"`
	Time      string `json:"time"`
	Oi        string `json:"oi"`         // 1=退出，2=加入
	GroupType int    `json:"group_type"` // 1=养生局; 2=竞技局; 3=强度局
}

type MsgBoard struct {
	LikeUsers  []string `json:"like_users"` // 点赞这条评价的所有用户
	GroupId    string   `json:"group_id" validate:"required"`
	User       string   `json:"user" validate:"required"` // 写下评价的用户或者是点赞这条评价的用户
	NickName   string   `json:"nick_name"`
	Img        string   `json:"img" validate:"required"`
	Evaluate   string   `json:"evaluate" validate:"required"`
	EvaluateId string   `json:"evaluate_id"`
	Time       string   `json:"time" validate:"required"`
	Like       int      `json:"like"`
	IsLike     bool     `json:"is_like"`
}

type UpdateVenueInfo struct {
	Id       string `json:"id" validate:"required"`        // 运动场地id
	Img      string `json:"img" validate:"required"`       // 场地图片
	City     string `json:"city" validate:"required"`      // 当前城市的中文名字：深圳市
	SportKey string `json:"sport_key" validate:"required"` // 运动类型: 篮球，足球....
	CityPy   string `json:"city_py"`                       // 城市的拼音名字: shenzhenshi
	Content  string `json:"content"`                       // 更新内容,目前只能统一更新图片,这里都写: 更新了场地图片
}

type FilterField struct {
	Id   int    `json:"id"`
	Type int    `json:"type"`
	Name string `json:"name"`
}

type OnlineData struct {
	Id       string `json:"id"`    // 场地id, 对应着每一个聊天室
	Title    string `json:"title"` // 场地简称，比如：科苑篮球馆
	City     string `json:"city"`
	SportKey string `json:"sport_key"` // 比如：篮球bks,完整的：shenzhenshi_bks
	Online   int    `json:"online"`    // 在线人数
}

type UserGetOnlineData struct {
	Id []string `json:"id" validate:"required"`
}

type WxBtnText struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type PublishData struct {
	UserCount   []*UserRoomID `json:"user_count"`                  // 沟通人数
	Id          string        `json:"id" validate:"required"`      // 生成唯一的任务id
	Rid         string        `json:"rid"`                         // room id
	City        string        `json:"city" validate:"required"`    // 城市：中文
	CityPy      string        `json:"city_py"`                     // 城市拼音：shenzhenshi
	UserId      string        `json:"user_id" validate:"required"` // 发布用户就用wx的openid
	NickName    string        `json:"nick_name" validate:"required"`
	Img         string        `json:"img" validate:"required"`
	Content     string        `json:"content" validate:"required"`
	Addr        string        `json:"addr" validate:"required"`
	Title       string        `json:"title" validate:"required"` // 篮球场简称
	Date        string        `json:"date" validate:"required"`  // 用户指定的陪练时间
	PublishDate string        `json:"publish_date"`              // 用户发布的时间，后端写
	Price       string        `json:"price" validate:"required"`
	GenderReq   string        `json:"gender_req" validate:"required"`
	SportKey    string        `json:"sport_key" validate:"required"`
	Time        string        `json:"time"` // 创建时间，后端写
	Players     string        `json:"players" validate:"required"`
	Lng         float64       `json:"lng" validate:"required"`
	Lat         float64       `json:"lat" validate:"required"`
	OnlineNum   int           `json:"online_num"`
	Finish      bool          `json:"finish"`
	IsDel       bool          `json:"is_del"`
	IsPublisher bool          `json:"is_publisher"`
}

type MissionStatus struct {
	Id       string `json:"id" validate:"required"` // 任务id
	UserId   string `json:"user_id" validate:"required"`
	City     string `json:"city" validate:"required"`
	CityPy   string `json:"city_py"`
	SportKey string `json:"sport_key" validate:"required"`
	Status   int    `json:"status" validate:"required"` // 1.表示完成，2.表示删除，3.撤销删除，4.撤销完成
}

type UserRoomID struct {
	Tid       string `json:"tid" validate:"required"`
	Rid       string `json:"rid"`
	NickName  string `json:"nick_name" validate:"required"`
	Img       string `json:"img" validate:"required"`
	UserId    string `json:"user_id" validate:"required"`
	City      string `json:"city" validate:"required"`
	UserCount int    `json:"user_count"`
}
