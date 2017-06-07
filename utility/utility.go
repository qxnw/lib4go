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

//EqualAndSet 判断变量的值与指定相等时设置为另一个值，否则使用原值
func EqualAndSet(def int, a int, b int) int {
	if def == a {
		return b
	}
	return def
}

//DecodeString 判断变量的值与指定相等时设置为另一个值，否则使用原值
func DecodeString(def string, a string, b string) string {
	if def == a {
		return b
	}
	return def
}

//DecodeInt 检查i是否等于a,相等则返回b,否则返回c
func DecodeInt(i interface{}, a int, b int, c int) int {
	if v, ok := i.(int); ok {
		if v == a {
			return b
		}
		return c
	}
	return c
}

//DecodeString2 检查i是否等于a,相等则返回b,否则返回c
func DecodeString2(i interface{}, a string, b string, c string) string {
	if v, ok := i.(string); ok {
		if v == a {
			return b
		}
		return c
	}
	return c
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
