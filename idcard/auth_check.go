package idcard

import (
	"fmt"
	"log"

	"github.com/cupen/game-anti-addiction/auth"
)

const (
	// idcard authentication
	url_check = "https://api.wlc.nppa.gov.cn/idcard/authentication/check"

	// idcard authentication for testcase
	url_check_test = "https://wlc.nppa.gov.cn/test/authentication/check/%s"
)

// CheckStatus ...
var CheckStatus = struct {
	Pass       int
	InProgress int
	Fail       int
}{
	Pass:       0,
	InProgress: 1,
	Fail:       2,
}

// Do ...
func (req *CheckRequest) Do(c *auth.Client, _url ...string) (*CheckResponse, error) {
	realUrl := url_check
	if len(_url) > 0 {
		realUrl = _url[0]
	}
	resp, err := c.PostJSON(realUrl, req)
	if err != nil {
		return nil, err
	}

	respObj := CheckResponse{ErrCode: -1}
	if err := resp.AsJson(&respObj); err != nil {
		log.Printf("err:%v %s,", err, resp.Dump())
		return nil, err
	}
	return &respObj, err
}

// DoTestSuite ...
func (req *CheckRequest) DoTestSuite(c *auth.Client, testCode string) (*CheckResponse, error) {
	_url := fmt.Sprintf(url_check_test, testCode)
	return req.Do(c, _url)
}

// IsOK ...
func (cr *CheckResponse) IsOK() bool {
	return cr.ErrCode == 0
}

// IsPassed ...
func (cr *CheckResponse) IsPassed() bool {
	if !cr.IsOK() {
		return false
	}
	return cr.Data.Result.Status == CheckStatus.Pass
}

// IsInProgress ...
func (cr *CheckResponse) IsInProgress() bool {
	if !cr.IsOK() {
		return false
	}
	return cr.Data.Result.Status == CheckStatus.InProgress
}
