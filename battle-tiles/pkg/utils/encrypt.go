package utils

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
)

// EncryptSHA256 使用 SHA256 加密并返回十六进制字符串
func EncryptSHA256(data string) string {
	hash := sha256.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}
func EncryptMD5(data string) string {
	hash := md5.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}
