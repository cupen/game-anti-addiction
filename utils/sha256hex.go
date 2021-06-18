package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

func SHA256ToHex(data []byte) (string, error) {
	h := sha256.New()
	_, err := h.Write(data)
	if err != nil {
		return "", err
	}
	rs := hex.EncodeToString(h.Sum(nil))
	// log.Printf("SHA256ToHex:\n\tsrc: %s \n\tdst: %s", string(data), rs)
	return rs, nil
}
