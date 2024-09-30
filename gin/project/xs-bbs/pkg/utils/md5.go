package utils

import (
	"crypto/md5"
	"encoding/hex"
)

// MD5 加密
func MD5String(s,secret string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(s)))
}
