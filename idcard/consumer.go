package idcard

import (
	"encoding/json"
	"time"

	"github.com/cupen/game-anti-addiction/auth"
	"go.uber.org/ratelimit"
)

func DecodeQueryRequest(msgList [][]byte) ([]*QueryRequest, error) {
	var rsList = []*QueryRequest{}
	var lastErr error = nil
	for _, msg := range msgList {
		obj := QueryRequest{}
		if err := json.Unmarshal(msg, &obj); err != nil {
			lastErr = err
			continue
		}
		rsList = append(rsList, &obj)
	}
	return rsList, lastErr
}

func ConsumerFunc(c *auth.Client, rate int, callback func(*QueryResponse)) func([][]byte) error {
	limiter := ratelimit.New(rate)
	return func(msgList [][]byte) error {
		if len(msgList) <= 0 {
			return nil
		}
		reqList, err := DecodeQueryRequest(msgList)
		if reqList == nil {
			return err
		}

		for _, req := range reqList {
			for i := 0; i < 3; i++ {
				if i > 0 {
					time.Sleep(time.Duration(i) * time.Second)
				}
				now := limiter.Take()
				_ = now
				resp, err := req.Do(c)
				if err != nil {
					continue
				}
				if callback != nil {
					callback(resp)
				}
				break
			}
		}
		return nil
	}
}
