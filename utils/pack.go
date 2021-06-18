package utils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

// Pack ...
func Pack(data []byte, key string) ([]byte, error) {
	cipherText, err := aes128gcm(data, key)
	if err != nil {
		return nil, err
	}
	packed := map[string]string{
		"data": base64.StdEncoding.EncodeToString(cipherText),
	}
	return json.Marshal(packed)
}

// PackForDebug ...
func PackForDebug(data []byte, key string, nonce string) ([]byte, error) {
	cipherText, err := aes128gcm_bynonce(data, key, nonce)
	if err != nil {
		return nil, err
	}
	packed := map[string]string{
		"data": base64.StdEncoding.EncodeToString(cipherText),
	}
	return json.Marshal(packed)
}

// UnPack ...
func UnPack(data []byte, key string) ([]byte, error) {
	packed := map[string]string{}
	if err := json.Unmarshal(data, &packed); err != nil {
		return nil, err
	}

	cipherTextB64, ok := packed["data"]
	if !ok {
		return nil, fmt.Errorf("invalid data-packed: missing field \"data\"")
	}
	cipherText, err := base64.StdEncoding.DecodeString(cipherTextB64)
	if err != nil {
		return nil, err
	}

	plainText, err := unaes128gcm(cipherText, key)
	if err != nil {
		return nil, err
	}
	return plainText, nil
}
