package logger

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func sysLoggerInfo(callBack func() string, content ...interface{}) {
	sysLoggerWrite(SLevel_Info, fmt.Sprint(content...), callBack)
}
func sysLoggerError(callBack func() string, content ...interface{}) {
	sysLoggerWrite(SLevel_Error, fmt.Sprint(content...), callBack)
}

func sysLoggerWrite(level string, content interface{}, callBack func() string) {
	if strings.EqualFold(level, "") {
		level = "All"
	}

	e := LogEvent{}
	e.Now = time.Now()
	e.Level = level
	e.Name = "sys"
	e.Session = getSessionID()
	e.Content = fmt.Sprintf("%v", content)
	e.Output = "[%datetime][%l][%session] %content%n"
	os.Stderr.WriteString(transform(e.Output, e))
	if callBack != nil {
		callBack()
	}
}
