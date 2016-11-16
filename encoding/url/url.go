package url

import (
	"net/url"
)

func Encode(input string) string {
	return url.QueryEscape(input)
}

func Decode(input string) (string, error) {
	return url.QueryUnescape(input)
}
