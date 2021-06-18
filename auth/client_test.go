package auth

import (
	"net/url"
	"testing"
)

func getTestCase1() (appId, bizId, secretKey, timestamp, body, sign string) {
	appId = "test-appId"
	bizId = "test-bizId"
	secretKey = "2836e95fcd10e04b0069bb1ee659955b"
	timestamp = "1584949895758"
	body = `{"data":"CqT/33f3jyoiYqT8MtxEFk3x2rlfhmgzhxpHqWosSj4d3hq2EbrtVyx2aLj565ZQNTcPrcDipnvpq/D/vQDaLKW70O83Q42zvR0//OfnYLcIjTPMnqa+SOhsjQrSdu66ySSORCAo"}`
	sign = "386c03b776a28c06b8032a958fbd89337424ef45c62d0422706cca633d8ad5fd"
	return
}

func TestClient_makeSign(t *testing.T) {
	appId, bizId, secretKey, timestamp, body, sign := getTestCase1()
	c := NewClient(appId, bizId, secretKey)
	qs := url.Values{}
	qs.Set("id", "test-id")
	qs.Set("name", "test-name")
	signActual, err := c.makeSign(timestamp, body, qs)
	if err != nil {
		t.Fatalf("%v", err)
		return
	}
	if sign != signActual {
		t.Fatalf("invalid sign.\nexpected: %s\nactual  : %s ", sign, signActual)
		return
	}
}
