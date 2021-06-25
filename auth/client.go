package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/cupen/game-anti-addiction/utils"
)

type Client struct {
	appId      string
	bizId      string
	secretKey  string
	httpClient *http.Client
	debug      bool
	debugArgs  url.Values
}

func NewClient(appId, bizId, secretKey string, opts ...Option) *Client {
	c := Client{
		appId:      appId,
		bizId:      bizId,
		secretKey:  secretKey,
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}
	return c.ApplyOptions(opts...)
}

func (c *Client) Clone() *Client {
	cloned := *c
	if c == &cloned {
		panic(fmt.Errorf("clone Client failed: %p", c))
	}
	return &cloned
}

func (c *Client) ApplyOptions(opts ...Option) *Client {
	for _, optFunc := range opts {
		optFunc(c)
	}
	return c
}

func (c *Client) NewHeaders(body string, now time.Time, qs url.Values) (*Headers, error) {
	ts := c.makeTimestamp(now)
	if c.IsDebug() {
		if timestamp := c.debugArgs.Get("timestamp"); timestamp != "" {
			ts = timestamp
			log.Printf("[debug] set timestamp = %s", timestamp)
		}
	}
	sign, err := c.makeSign(ts, body, qs)
	if err != nil {
		return nil, err
	}
	return &Headers{
		AppID:     c.appId,
		BizID:     c.bizId,
		Timestamp: ts,
		Sign:      sign,
	}, nil
}

func (c *Client) PostJSON(url string, body interface{}) (*ReqRespContext, error) {
	var _body []byte
	var err error
	_body, ok := body.([]byte)
	if !ok {
		_body, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}

	bodyPacked, err := utils.Pack(_body, c.secretKey)
	if c.IsDebug() {
		if nonce := c.debugArgs.Get("nonce"); nonce != "" {
			log.Printf("[debug] set nonce = %s", nonce)
			bodyPacked, err = utils.PackForDebug(_body, c.secretKey, nonce)
		}
	}
	if err != nil {
		return nil, err
	}

	headers, err := c.NewHeaders(string(bodyPacked), time.Now(), nil)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewReader(bodyPacked)
	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return nil, err
	}

	headers.SetHTTPHeaders(req.Header)
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	rs := &ReqRespContext{
		Method:     "POST-JSON",
		URL:        url,
		headers:    req.Header,
		reqDataRaw: _body,
		reqData:    bodyPacked,
		StatusCode: resp.StatusCode,
		respBody:   respData,
	}
	if c.IsDebug() {
		log.Println(rs.Dump())
	}
	return rs, nil
}

func (c *Client) Get(url string, qs url.Values) (*ReqRespContext, error) {
	if qs != nil {
		url += ("?" + qs.Encode())
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	headers, err := c.NewHeaders("", time.Now(), qs)
	if err != nil {
		return nil, err
	}
	headers.SetHTTPHeaders(req.Header)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	rs := &ReqRespContext{
		Method:     "GET",
		URL:        url,
		headers:    req.Header,
		reqDataRaw: nil,
		reqData:    nil,
		StatusCode: resp.StatusCode,
		respBody:   respData,
	}
	if c.IsDebug() {
		log.Println(rs.Dump())
	}
	return rs, nil
}

func (c *Client) makeTimestamp(now time.Time) string {
	ts := now.UnixNano() / int64(time.Millisecond)
	return strconv.FormatInt(ts, 10)
}

func (c *Client) makeSign(ts, body string, qs url.Values) (string, error) {
	return c.newSign(ts, body, qs)
}

func (c *Client) newSign(ts, body string, qs url.Values) (string, error) {
	params := map[string]string{}
	for k, v := range qs {
		if len(v) <= 0 {
			continue
		}
		params[k] = v[0]
	}
	params["appId"] = c.appId
	params["bizId"] = c.bizId
	params["timestamps"] = ts

	keys := generateSortedKeys(params)
	buf := bytes.NewBuffer([]byte(c.secretKey))
	for _, k := range keys {
		buf.WriteString(k)
		buf.WriteString(params[k])
	}
	if body != "" {
		buf.WriteString(body)
	}
	sign, err := utils.SHA256ToHex(buf.Bytes())
	if err != nil {
		return "", err
	}

	if c.IsDebug() {
		log.Printf("[debug] sign:\n\tsrt=%s\n\tdst=%s", buf.String(), sign)
	}

	if sign == "" {
		return "", fmt.Errorf("empty sign: unknow error")
	}
	return sign, nil
}

func (c *Client) IsDebug() bool {
	return c.debug
}
