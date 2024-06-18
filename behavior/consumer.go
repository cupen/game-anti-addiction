package behavior

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/cupen/game-anti-addiction/auth"
	"go.uber.org/ratelimit"
)

func DecodeLoginOutRequest(msgList [][]byte, batchSize int) ([]*LoginOutRequest, error) {
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
			rsList = append(rsList, &LoginOutRequest{Collections: events})
			events = []LoginOutEvent{}
		}
	}
	if len(events) > 0 {
		rsList = append(rsList, &LoginOutRequest{Collections: events})
	}
	return rsList, lastErr
}

func ConsumerFunc(c *auth.Client, batchSize, rate int) func([][]byte) error {
	limiter := ratelimit.New(rate)
	return func(msgList [][]byte) error {
		if len(msgList) <= 0 {
			return nil
		}
		reqList, err := DecodeLoginOutRequest(msgList, batchSize)
		if reqList == nil {
			return err
		}
		var lastErr error
		for _, req := range reqList {
			for i := 0; i < 3; i++ {
				if i > 0 {
					time.Sleep(time.Duration(i) * time.Second)
				}
				_ = limiter.Take()
				resp, err := req.Do(c)
				if resp == nil {
					lastErr = fmt.Errorf("nil response: %v", err)
					log.Printf("[GAA] push LoginOutEvent failed: %v", lastErr)
					continue
				}
				lastErr = err
				if resp.CanRetry() {
					log.Printf("[GAA] push LoginOutEvent retrying")
					continue
				}
				if !resp.IsOK() {
					lastErr = resp.AsError()
					log.Printf("[GAA] push LoginOutEvent failed: %v", lastErr)
				} else {
					lastErr = nil
				}
				break
			}
		}
		return lastErr
	}
}
