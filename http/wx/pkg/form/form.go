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
	JoinUsers              []*JoinGroupUsers `json:"join_users"` // 某个运动场地，用户点击加入组件的人数
	UserReviews            []*MsgBoard       `json:"user_reviews"`
	VenueUpdateUsers       []*AddrListForm   `json:"venue_update_users"`
	Tags                   []string          `json:"tags"`
	Images                 []string          `json:"images"`
	Id                     string            `json:"id"`
	Img                    string            `json:"img"`
	Addr                   string            `json:"addr"`
	Title                  string            `json:"title"`
	UserId                 string            `json:"user_id"`
	Online                 string            `json:"online"`
	Distance               string            `json:"distance"`
	Aid                    string            `json:"aid"` // 接口返回的地址唯一id，再次请求接口返回的id是一致的，更新的时候有用
	JoinUserCount          int               `json:"join_user_count"`
	UserReviewsCount       int               `json:"user_reviews_count"`
	VenueUpdateUsersCount  int               `json:"venue_update_users_count"`
	Lng                    float64           `json:"lng"`
	Lat                    float64           `json:"lat"`
	DisVal                 float64           `json:"dis_val"`
	IsShow                 bool              `json:"is_show"`
	IsShowUserReviews      bool              `json:"is_show_user_reviews"`
	IsShowJoinUsers        bool              `json:"is_show_join_users"`
	IsShowVenueUpdateUsers bool              `json:"is_show_venue_update_users"`
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
	Name    string `json:"name"`
	Key     string `json:"key"`
	Icon    string `json:"icon"`
	Img     string `json:"img"`
	Checked bool   `json:"checked"`
}

type GroupOnlineStatus struct {
	GroupId    string `json:"group_id"`
	OnlineUser string `json:"online_user"`
}

type JoinGroupUsers struct {
	GroupId  string `json:"group_id" validate:"required"`
	User     string `json:"user"`
	NickName string `json:"nick_name"`
	Img      string `json:"img"`
	Style    string `json:"style"`
	Skill    string `json:"skill"`
	Time     string `json:"time"`
	Oi       string `json:"oi"`
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
