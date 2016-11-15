package md5

import (
	"crypto/md5"
	"encoding/hex"
)

func Encrypt(s string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(s))
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}
