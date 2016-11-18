package url

import (
	"net/url"
)

// Encode 对字符串进行url编码
func Encode(input string) string {
	return url.QueryEscape(input)
}

// Decode 对字符串进行url解码
func Decode(input string) (string, error) {
	return url.QueryUnescape(input)
}
