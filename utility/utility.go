package utility

import (
	"crypto/rand"
	"encoding/base64"
	"io"

	"github.com/qxnw/lib4go/security/md5"
)

// GetGUID 生成Guid字串
func GetGUID() string {
	b := make([]byte, 48)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return md5.Encrypt(base64.URLEncoding.EncodeToString(b))
}

//IsStringEmpty 当前对像是否是字符串空
func IsStringEmpty(v interface{}) bool {
	if v == nil {
		return true
	}
	if t, ok := v.(string); ok && len(t) == 0 {
		return true
	}
	return false
}
