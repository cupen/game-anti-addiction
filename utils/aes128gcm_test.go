package utils

import (
	"encoding/base64"
	"encoding/json"
	"testing"
)

func TestAes128gcm(t *testing.T) {
	secretKey, plainText, _ := testData1()
	cipherText, err := aes128gcm([]byte(plainText), secretKey)
	if err != nil {
		t.Fatalf("%v", err)
		return
	}
	plainTextActual, err := unaes128gcm(cipherText, secretKey)
	if err != nil {
		t.Fatalf("%v", err)
		return
	}
	if plainText != string(plainTextActual) {
		t.Fatalf("invalid result of aes128gcm.\nexpected:%s\nactual  :%s", plainText, string(plainTextActual))
	}
}

func TestAes128gcm_bynonce(t *testing.T) {
	secretKey, plainText, packedText := testData1()
	m := map[string]string{}
	_ = json.Unmarshal([]byte(packedText), &m)
	cipherText := m["data"]
	nonce := cipherText[0:16]
	_cipherTextActual, err := aes128gcm_bynonce([]byte(plainText), secretKey, nonce)
	if err != nil {
		t.Fatalf("%v", err)
		return
	}
	cipherTextActual := base64.StdEncoding.EncodeToString(_cipherTextActual)
	if cipherText != string(cipherTextActual) {
		t.Fatalf("invalid result of aes128gcm.\nexpected:%s\nactual  :%s", cipherText, string(cipherTextActual))
	}
}

func TestUnAes128gcm(t *testing.T) {
	secretKey, plainText, packedText := testData1()
	m := map[string]string{}
	_ = json.Unmarshal([]byte(packedText), &m)
	cipherText := m["data"]

	cipherTextBytes, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		t.Fatalf("%v", err)
		return
	}

	plainTextActual, err := unaes128gcm(cipherTextBytes, secretKey)
	if err != nil {
		t.Fatalf("%v", err)
		return
	}
	if plainText != string(plainTextActual) {
		t.Fatalf("invalid result of pack.\nexpected:%s\nactual  :%s", plainText, string(plainTextActual))
		return
	}
}
