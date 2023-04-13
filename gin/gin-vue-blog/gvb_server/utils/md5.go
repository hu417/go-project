package utils

import (
	"crypto/md5"
	"encoding/hex"
)

// md5
func Md5(data []byte) string {
	m := md5.New()
	m.Write(data)
	res := hex.EncodeToString(m.Sum(nil))
	return res
}
