package utils

import (
	"testing"
)

func TestSHA256ToHex(t *testing.T) {
	src := `2836e95fcd10e04b0069bb1ee659955bappIdtest-appIdbizIdtest-bizIdidtest-idnametest-nametimestamps1584949895758{"data":"CqT/33f3jyoiYqT8MtxEFk3x2rlfhmgzhxpHqWosSj4d3hq2EbrtVyx2aLj565ZQNTcPrcDipnvpq/D/vQDaLKW70O83Q42zvR0//OfnYLcIjTPMnqa+SOhsjQrSdu66ySSORCAo"}`
	dst, err := SHA256ToHex([]byte(src))
	if err != nil {
		t.Fatalf("invalid sha256:%v", err)
	}

	expected := "386c03b776a28c06b8032a958fbd89337424ef45c62d0422706cca633d8ad5fd"
	if dst != expected {
		t.Fatalf("invalid sha256hex.\nexpected:%s\nactual  :%s", expected, dst)
	}
}
