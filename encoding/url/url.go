package url

import (
	"net/url"
)

// 对字符串进行url编码
func Encode(input string) string {
	return url.QueryEscape(input)
}

// 对字符串进行url解码
func Decode(input string) (string, error) {
	return url.QueryUnescape(input)
}
