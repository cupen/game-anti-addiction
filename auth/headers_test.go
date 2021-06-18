package auth

import (
	"reflect"
	"testing"
)

func Test_generateSortedKeys(t *testing.T) {
	obj := Headers{
		AppID:     "value(appId)",
		BizID:     "value(bizId)",
		Timestamp: "value(ts)",
	}
	params := map[string]string{
		"appId":      obj.AppID,
		"bizId":      obj.BizID,
		"timestamps": obj.Timestamp,
		"ai":         "value(ai)",
	}
	keys := generateSortedKeys(params)
	keysExpected := []string{
		"ai",
		"appId",
		"bizId",
		"timestamps",
	}
	if !reflect.DeepEqual(keysExpected, keys) {
		t.Fatalf("invalid keys:%v", keys)
	}
}
