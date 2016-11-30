package json

import "strings"

// Escape 把编码 \\u0026，\\u003c，\\u003e 替换为 &,<,>
func Escape(input string) string {
	r := strings.Replace(input, "\\u0026", "&", -1)
	r = strings.Replace(r, "\\u003c", "<", -1)
	r = strings.Replace(r, "\\u003e", ">", -1)
	return r
}
