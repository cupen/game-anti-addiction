package behavior

var (

	// 行为类型
	BehaviorTypes = struct {
		// 上线
		Online int
		// 下线
		Offline int
	}{
		Online:  0,
		Offline: 1,
	}

	// 用户类型
	UserTypes = struct {
		// 0 - 认证用户
		UserAuthed int
		// 2 - 游客
		Guest int
	}{
		UserAuthed: 0,
		Guest:      2,
	}
)

type LoginOutRequest struct {
	Collections []LoginOutEvent `json:"collections"`
}

type LoginOutResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	Data    struct {
		Results []LoginOutResult `json:"results"`
	} `json:"data"`
}

func (lor *LoginOutResponse) IsOK() bool {
	return lor.ErrCode == 0
}

func (lor *LoginOutResponse) CanRetry() bool {
	return lor.ErrCode == 1005 || lor.ErrCode == 1006
}

type LoginOutEvent struct {
	Num          int    `json:"no"`
	SessionID    string `json:"si"`
	BehaviorType int    `json:"bt"`
	Timestamp    int64  `json:"ot"`
	UserType     int    `json:"ct"`
	DeviceID     string `json:"di"`
	PlayerID     string `json:"pi"`
}

type LoginOutResult struct {
	Num     int    `json:"no"`
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}
