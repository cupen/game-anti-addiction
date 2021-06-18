package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
)

// aes128gcm ...
func aes128gcm(plainText []byte, key string) ([]byte, error) {
	_key, err := hex.DecodeString(key)
	if err != nil {
		return nil, err
	}
	c, err := aes.NewCipher(_key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	dst := nonce
	rs := gcm.Seal(dst, nonce, plainText, nil)
	return rs, nil
}

// aes128gcm_bynonce ...
func aes128gcm_bynonce(plantText []byte, keyHex string, nonceb64 string) ([]byte, error) {
	nonce, err := base64.StdEncoding.DecodeString(nonceb64)
	if err != nil {
		return nil, err
	}
	_key, err := hex.DecodeString(keyHex)
	if err != nil {
		return nil, err
	}
	c, err := aes.NewCipher(_key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}
	// nonce := make([]byte, gcm.NonceSize())
	// if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
	// 	return nil, err
	// }
	if len(nonce) != gcm.NonceSize() {
		return nil, fmt.Errorf("invalid nonce, it must be bytes.length=%d. b64:%s", len(nonce), nonceb64)
	}
	dst := nonce
	rs := gcm.Seal(dst, nonce, plantText, nil)
	return rs, nil
}

// unaes128gcm ...
func unaes128gcm(cipherText []byte, key string) ([]byte, error) {
	_key, err := hex.DecodeString(key)
	if err != nil {
		return nil, err
	}
	c, err := aes.NewCipher(_key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(cipherText) < nonceSize {
		return nil, fmt.Errorf("invalid nonce size. %d < %d", len(cipherText), nonceSize)
	}

	nonce := cipherText[:nonceSize]
	cipherText = cipherText[nonceSize:]
	plainText, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return nil, err
	}
	return plainText, nil
}
