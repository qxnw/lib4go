package logger

import (
	"fmt"
	"os"
	"time"
)

func sysLoggerInfo(content ...interface{}) {
	sysLoggerWrite(SLevel_Info, fmt.Sprint(content...))
}
func sysLoggerError(content ...interface{}) {
	sysLoggerWrite(SLevel_Error, fmt.Sprint(content...))
}

func sysLoggerWrite(level string, content interface{}) {
	e := LogEvent{}
	e.Now = time.Now()
	e.Level = level
	e.Name = "sys"
	e.Session = getSessionID()
	e.Content = fmt.Sprintf("%v", content)
	e.Output = "[%datetime][%l][%session] %content%n"
	os.Stderr.WriteString(transform(e.Output, e))
}
