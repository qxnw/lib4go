package url

import (
	"html"
	"net/url"
)

func URLEncode(input string) string {
	return url.QueryEscape(input)
}

func URLDecode(input string) (string, error) {
	return url.QueryUnescape(input)
}

func HTMLEncode(input string) string {
	return html.EscapeString(input)
}

func HTMLDecode(input string) string {
	return html.UnescapeString(input)
}
