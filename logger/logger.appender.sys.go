package logger

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func sysLoggerInfo(callBack func(error), content ...interface{}) {
	sysLoggerWrite(callBack, SLevel_Info, fmt.Sprint(content...))
}
func sysLoggerError(callBack func(error), content ...interface{}) {
	sysLoggerWrite(callBack, SLevel_Error, fmt.Sprint(content...))
}

func sysLoggerWrite(callBack func(error), level string, content interface{}) {
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
		callBack(content.(error))
	}
}
