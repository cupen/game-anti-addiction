package utils

import "testing"

// 加密示例
// //访问密钥
// secretKey:2836e95fcd10e04b0069bb1ee659955b
// //待加密数据
// {"ai":"test-accountId","name":"用户姓名","idNum":"371321199012310912"}
// //加密后请求体数据
// {"data":"CqT/33f3jyoiYqT8MtxEFk3x2rlfhmgzhxpHqWosSj4d3hq2EbrtVyx2aLj5
// 65ZQNTcPrcDipnvpq/D/vQDaLKW70O83Q42zvR0//OfnYLcIjTPMnqa+SOhsjQrSdu66y
// SSORCAo"}
func testData1() (secretKey, plainText, packedText string) {
	secretKey = "2836e95fcd10e04b0069bb1ee659955b"
	plainText = `{"ai":"test-accountId","name":"用户姓名","idNum":"371321199012310912"}`
	packedText = `{"data":"CqT/33f3jyoiYqT8MtxEFk3x2rlfhmgzhxpHqWosSj4d3hq2EbrtVyx2aLj565ZQNTcPrcDipnvpq/D/vQDaLKW70O83Q42zvR0//OfnYLcIjTPMnqa+SOhsjQrSdu66ySSORCAo"}`
	return
}

func TestPack(t *testing.T) {
	secretKey, plainText, _ := testData1()
	packedTextActual, err := Pack([]byte(plainText), secretKey)
	if err != nil {
		t.Fatalf("%v", err)
		return
	}
	plainTextActual, err := UnPack([]byte(packedTextActual), secretKey)
	if err != nil {
		t.Fatalf("%v", err)
		return
	}
	if plainText != string(plainTextActual) || len(plainTextActual) <= 0 {
		t.Fatalf("invalid result of pack.\nexpected:%s\nactual  :%s", plainText, string(plainTextActual))
		return
	}
}

func TestUnPack(t *testing.T) {
	secretKey, plainText, packedText := testData1()
	plainTextActual, err := UnPack([]byte(packedText), secretKey)
	if err != nil {
		t.Fatalf("%v", err)
		return
	}
	if plainText != string(plainTextActual) {
		t.Fatalf("invalid result of pack.\nexpected:%s\nactual  :%s", plainText, string(plainTextActual))
		return
	}
}
