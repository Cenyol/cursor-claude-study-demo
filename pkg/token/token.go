package token

import (
	"crypto/rand"
	"encoding/hex"
)

// New 生成随机 Token（Session 用）
func New() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
