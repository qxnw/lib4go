package html

import (
	"html"
)

func HTMLEncode(input string) string {
	return html.EscapeString(input)
}

func HTMLDecode(input string) string {
	return html.UnescapeString(input)
}
