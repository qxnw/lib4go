package logger

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

func transform(tpl string, event LogEvent) (result string) {
	word, _ := regexp.Compile(`%\w+`)
	//@变量, 将数据放入params中
	result = word.ReplaceAllStringFunc(tpl, func(s string) string {
		key := s[1:]
		switch key {
		case "session":
			return event.Session
		case "date":
			return event.Now.Format("20060102")
		case "datetime":
			return event.Now.Format("2006/01/02 15:04:05")
		case "year":
			return event.Now.Format("2006")
		case "mm":
			return event.Now.Format("01")
		case "dd":
			return event.Now.Format("02")
		case "hh":
			return event.Now.Format("15")
		case "mi":
			return event.Now.Format("04")
		case "ss":
			return event.Now.Format("05")
		case "level":
			return strings.ToLower(event.Level)
		case "l":
			return strings.ToLower(event.Level)[:1]
		case "name":
			return event.Name
		case "pid":
			return fmt.Sprintf("%d", os.Getpid())
		case "n":
			return "\n"
		case "content":
			return event.Content
		default:
			v, ok := event.Tags.Get(key)
			if ok {
				return v.(string)
			}
			return ""
		}
	})
	return
}
