package behavior

import (
	"encoding/json"
	"time"

	"github.com/cupen/game-anti-addiction/auth"
	"go.uber.org/ratelimit"
)

func decodeLoginoutRequest(msgList [][]byte, batchSize int) ([]*LoginOutRequest, error) {
	var rsList = []*LoginOutRequest{}
	var events = []LoginOutEvent{}
	var i = 0
	var lastErr error = nil
	for _, msg := range msgList {
		obj := LoginOutEvent{}
		if err := json.Unmarshal(msg, &obj); err != nil {
			lastErr = err
			continue
		}
		i %= batchSize
		i++
		obj.Num = i
		events = append(events, obj)
		if len(events) >= batchSize {
			req := LoginOutRequest{Collections: events}
			rsList = append(rsList, &req)
		}
	}
	if len(events) <= 0 {
		return nil, lastErr
	}
	return rsList, lastErr
}

func Consumer(c *auth.Client, batchSize, rate int) func([][]byte) error {
	limiter := ratelimit.New(rate)

	return func(msgList [][]byte) error {
		if len(msgList) <= 0 {
			return nil
		}
		reqList, err := decodeLoginoutRequest(msgList, batchSize)
		if reqList == nil {
			return err
		}

		for _, req := range reqList {
			for i := 0; i < 3; i++ {
				if i > 0 {
					time.Sleep(time.Duration(i) * time.Second)
				}
				_ = limiter.Take()
				resp, err := req.Do(c)
				if err != nil {
					continue
				}

				if resp.ErrCode != 0 && resp.CanRetry() {
					continue
				}
				break
			}
		}
		return nil
	}
}
