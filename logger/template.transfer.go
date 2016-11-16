package logger

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

func transform(format string, event LogEvent) string {
	var resultString string
	resultString = format
	formater := make(map[string]string)
	formater["session"] = event.Session
	formater["date"] = event.Now.Format("20060102")
	formater["datetime"] = event.Now.Format("2016/01/02 15:04:05")
	formater["year"] = event.Now.Format("2006")
	formater["mm"] = event.Now.Format("01")
	formater["dd"] = event.Now.Format("02")
	formater["hh"] = event.Now.Format("15")
	formater["mi"] = event.Now.Format("04")
	formater["ss"] = event.Now.Format("05")
	formater["level"] = strings.ToLower(event.Level)
	formater["l"] = strings.ToLower(event.Level)[:1]
	formater["name"] = event.Name
	formater["pid"] = fmt.Sprintf("%d", os.Getpid())
	formater["msg"] = event.Content
	for i, v := range formater {
		match, _ := regexp.Compile("%" + i)
		resultString = match.ReplaceAllString(resultString, v)
	}
	return resultString
}
