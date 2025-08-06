package openid

var (
	wl = []string{"ogR3E62jXXJMbVcImRqMA1gTSegM", "user_ogR3E62jXXJMbVcImRqMA1gTSegM", "user_ogR3E6zO5aFz3RF-CzI14q7U_sHI", "ogR3E6zO5aFz3RF-CzI14q7U_sHI"}
)

type WhiteList struct {
	uid string
}

func NewWhiteList(uid string) WhiteList {
	return WhiteList{
		uid: uid,
	}
}

func (w WhiteList) IsWhite() bool {
	for _, u := range wl {
		if u == w.uid {
			return true
		}
	}

	return false
}
