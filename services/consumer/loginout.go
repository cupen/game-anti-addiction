package consumer

import (
	"encoding/json"
	"time"

	"github.com/cupen/game-anti-addiction/auth"
	"github.com/cupen/game-anti-addiction/behavior"
)

func NewLoginoutRequest(msgList [][]byte) (*behavior.LoginOutRequest, error) {
	var objList = []behavior.LoginOutEvent{}
	var i = 0
	var lastErr error = nil
	for _, msg := range msgList {
		obj := behavior.LoginOutEvent{}
		if err := json.Unmarshal(msg, &obj); err != nil {
			lastErr = err
			continue
		}
		i++
		obj.Num = i
		objList = append(objList, obj)
	}
	if len(objList) <= 0 {
		return nil, lastErr
	}
	return &behavior.LoginOutRequest{Collections: objList}, lastErr
}

func HandlerForLoginoutRequest(c *auth.Client) Handler {
	return func(msgList [][]byte) error {
		if len(msgList) <= 0 {
			return nil
		}
		req, err := NewLoginoutRequest(msgList)
		if req == nil {
			return err
		}

		for i := 0; i < 10; i++ {
			if i > 0 {
				time.Sleep(time.Duration(i) * time.Second)
			}

			resp, err := req.Do(c)
			if err != nil {
				continue
			}
			if resp.ErrCode == 0 {
				break
			}
		}
		return nil
	}
}
