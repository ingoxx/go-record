package form

type AddAddrForm struct {
	Id     string  `json:"id" validate:"required"`
	Addr   string  `json:"addr" validate:"required"`
	UserId string  `json:"user_id" validate:"required"`
	Lat    float64 `json:"lat"`
	Lng    float64 `json:"Lng"`
}

type PassAddrReqForm struct {
	Id string `json:"id" validate:"required"`
}
