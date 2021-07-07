package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type ReqRespContext struct {
	Method     string
	URL        string
	StatusCode int
	reqDataRaw []byte
	reqData    []byte
	respBody   []byte
	headers    http.Header
}

func (r *ReqRespContext) AsJson(obj interface{}) error {
	return json.Unmarshal(r.respBody, obj)
}

func (r *ReqRespContext) Bytes() []byte {
	return r.respBody
}

func (r *ReqRespContext) Text() string {
	return string(r.respBody)
}

func (r *ReqRespContext) Dump() string {
	if r == nil {
		return "nil"
	}
	headers := ""
	if r.headers != nil {
		_data, _ := json.Marshal(r.headers)
		headers = string(_data)
	}
	lines := []string{
		r.Method,
		"\n\turl    : ", r.URL,
		"\n\theaders: ", headers,
		"\n\treq    : ", string(r.reqData),
		"\n\treq-raw: ", string(r.reqDataRaw),
		fmt.Sprintf("\n\tresp[%d] :", r.StatusCode),
		r.Text(),
	}
	return strings.Join(lines, "")
}
