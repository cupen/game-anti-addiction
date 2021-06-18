package idcard

// CheckRequest ...
type CheckRequest struct {
	AI    string `json:"ai"`
	Name  string `json:"name"`
	IDNum string `json:"idNum"`
}

// CheckResponse ...
type CheckResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	Data    struct {
		Result struct {
			Status int    `json:"status"`
			PI     string `json:"pi"`
		} `json:"result"`
	} `json:"data"`
}

// QueryRequest ...
type QueryRequest struct {
	AI string `json:"ai"`
}

// QueryResponse ...
type QueryResponse struct {
	CheckResponse
}
