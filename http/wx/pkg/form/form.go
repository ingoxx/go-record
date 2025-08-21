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
	JoinUsers        []JoinGroupUsers `json:"join_users"` // 某个运动场地，用户点击加入组件的人数
	UserReviews      []MsgBoard       `json:"user_reviews"`
	Tags             []string         `json:"tags"`
	Id               string           `json:"id"`
	Img              string           `json:"img"`
	Images           []string         `json:"images"`
	Addr             string           `json:"addr"`
	Title            string           `json:"title"`
	UserId           string           `json:"user_id"`
	Online           string           `json:"online"`
	Distance         string           `json:"distance"`
	Aid              string           `json:"aid"` // 接口返回的地址唯一id，再次请求接口返回的id是一致的，更新的时候有用
	JoinUserCount    int              `json:"join_user_count"`
	UserReviewsCount int              `json:"user_reviews_count"`
	Lng              float64          `json:"lng"`
	Lat              float64          `json:"lat"`
	IsShow           bool             `json:"is_show"`
}

type AddrListForm struct {
	Tags       string  `json:"tags"  validate:"required"`
	Id         string  `json:"id" validate:"required"`
	Addr       string  `json:"addr" validate:"required"`
	UserId     string  `json:"user_id" validate:"required"`
	City       string  `json:"city"  validate:"required"`        // 前端传入的是中文
	CityPy     string  `json:"city_py"`                          // 前端传入的中文转成拼音
	SportKey   string  `json:"sport_key" validate:"required"`    // 运动分类
	UpdateType string  `json:"update_type"  validate:"required"` // 更新类型：1.用户添加的新场地，2.用户更新了场地
	Aid        string  `json:"aid"`                              // 场地的唯一id
	Img        string  `json:"img"`
	Lat        float64 `json:"lat"`
	Lng        float64 `json:"lng"`
	IsRecord   bool    `json:"is_record"` // true：已记录（审核通过），false：未记录（还未审核通过）
	IsShow     bool    `json:"is_show"`   // 审核列表中的数据，true：隐藏，false：不隐藏

}

type PassAddrReqForm struct {
	Id         string `json:"id" validate:"required"`
	City       string `json:"city"  validate:"required"`
	UpdateType string `json:"update_type"  validate:"required"` // 更新类型：1.用户添加的新场地，2.用户更新了场地
	Img        string `json:"img" validate:"required"`
}

type WxOpenidList struct {
	Openid   string `json:"openid"`
	Img      string `json:"img"`
	NickName string `json:"nick_name"`
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
	GroupId string `json:"group_id" validate:"required"`
	User    string `json:"user"`
	Img     string `json:"img"`
	Oi      string `json:"oi"`
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
	Id       string `json:"id" validate:"required"`
	Img      string `json:"img" validate:"required"`
	City     string `json:"city" validate:"required"`      // 当前城市的中文名字
	SportKey string `json:"sport_key" validate:"required"` // 运动类型: 篮球，足球....
	CityPy   string `json:"city_py"`                       // 当前城市的拼音名字

}
