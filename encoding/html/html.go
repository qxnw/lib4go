package html

import (
	"html"
)

// html编码
func Encode(input string) string {
	return html.EscapeString(input)
}

// html解码
func Decode(input string) string {
	return html.UnescapeString(input)
}
