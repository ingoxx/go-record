package form

type AddrListForm struct {
	Tags     string  `json:"tags"  validate:"required"`
	Id       string  `json:"id" validate:"required"`
	Addr     string  `json:"addr" validate:"required"`
	UserId   string  `json:"user_id" validate:"required"`
	City     string  `json:"city"  validate:"required"`     // 前端传入的是中文
	CityPy   string  `json:"city_py"`                       // 前端传入的中文转成拼音
	SportKey string  `json:"sport_key" validate:"required"` // 运动分类
	Img      string  `json:"img"`
	Lat      float64 `json:"lat"`
	Lng      float64 `json:"lng"`
	IsRecord bool    `json:"is_record"` // true：已记录，false：未记录
	IsShow   bool    `json:"is_show"`   // 审核列表中的数据，true：隐藏，false：不隐藏
}

type PassAddrReqForm struct {
	Id   string `json:"id" validate:"required"`
	City string `json:"city"  validate:"required"`
}

type WxOpenidList struct {
	Openid string `json:"openid"`
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
