package behavior

import (
	"fmt"
	"log"
	"time"

	"github.com/cupen/game-anti-addiction/auth"
)

const (
	url_loginout      = "http://api2.wlc.nppa.gov.cn/behavior/collection/loginout"
	url_loginout_test = "https://wlc.nppa.gov.cn/test/collection/loginout/%s"
)

func (req *LoginOutRequest) Do(c *auth.Client, _url ...string) (*LoginOutResponse, error) {
	realUrl := url_loginout
	if len(_url) > 0 {
		realUrl = _url[0]
	}
	resp, err := c.PostJSON(realUrl, req)
	if err != nil {
		return nil, err
	}

	respObj := LoginOutResponse{}
	if err := resp.AsJson(&respObj); err != nil {
		log.Printf("err:%v %s", err, resp.Dump())
		return nil, err
	}
	return &respObj, nil
}

func (req *LoginOutRequest) DoTestSuite(c *auth.Client, testCode string) (*LoginOutResponse, error) {
	_url := fmt.Sprintf(url_loginout_test, testCode)
	return req.Do(c, _url)
}

func NewLoginEvent(si string, ts time.Time, userType int, deviceId, playerId string) LoginOutEvent {
	if userType == UserTypes.Guest {
		playerId = ""
	}
	return LoginOutEvent{
		SessionID:    si,
		BehaviorType: BehaviorTypes.Online,
		Timestamp:    ts.UnixNano() / int64(time.Millisecond),
		UserType:     userType,
		DeviceID:     deviceId,
		PlayerID:     playerId,
	}
}

func NewLogoutEvent(si string, ts time.Time, userType int, deviceId, playerId string) LoginOutEvent {
	if userType == UserTypes.Guest {
		playerId = ""
	}
	return LoginOutEvent{
		SessionID:    si,
		BehaviorType: BehaviorTypes.Offline,
		Timestamp:    ts.UnixNano() / int64(time.Millisecond),
		UserType:     userType,
		DeviceID:     deviceId,
		PlayerID:     playerId,
	}
}
