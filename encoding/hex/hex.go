package hex

import (
	"encoding/hex"
)

// 把[]byte类型通过hex编码成string
func Encode(src []byte) string {
	return hex.EncodeToString(src)
}

// 把一个string类型通过hex解码成string
func Decode(src string) (r string, err error) {
	data, err := hex.DecodeString(src)
	if err != nil {
		return
	}
	r = string(data)
	return
}
