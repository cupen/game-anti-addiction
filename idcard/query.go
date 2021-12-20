package idcard

import (
	"fmt"
	"log"
	"net/url"

	"github.com/cupen/game-anti-addiction/auth"
)

const (
	// idcard authentication
	url_query      = "https://api.wlc.nppa.gov.cn/idcard/authentication/query"
	url_query_test = "https://wlc.nppa.gov.cn/test/authentication/query/%s"
)

func (req *QueryRequest) Do(c *auth.Client, _url ...string) (*QueryResponse, error) {
	realUrl := url_query
	if len(_url) > 0 {
		realUrl = _url[0]
	}
	qs := url.Values{}
	qs.Set("ai", req.AI)
	resp, err := c.Get(realUrl, qs)
	if err != nil {
		return nil, err
	}

	respObj := QueryResponse{}
	respObj.ErrCode = -1
	if err := resp.AsJson(&respObj); err != nil {
		log.Printf("err:%v %s,", err, resp.Dump())
		return nil, err
	}
	return &respObj, err
}

func (req *QueryRequest) DoTestSuite(c *auth.Client, testCode string) (*QueryResponse, error) {
	_url := fmt.Sprintf(url_query_test, testCode)
	return req.Do(c, _url)
}
