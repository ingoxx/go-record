package form

type AddAddrForm struct {
	Id       string  `json:"id" validate:"required"`
	Addr     string  `json:"addr" validate:"required"`
	UserId   string  `json:"user_id" validate:"required"`
	City     string  `json:"city"  validate:"required"` // 前端传入的是中文
	CityPy   string  `json:"city_py"`                   // 前端传入的中文转成拼音
	Lat      float64 `json:"lat"`
	Lng      float64 `json:"lng"`
	IsRecord bool    `json:"is_record"`
}

type PassAddrReqForm struct {
	Id   string `json:"id" validate:"required"`
	City string `json:"city"  validate:"required"`
}
