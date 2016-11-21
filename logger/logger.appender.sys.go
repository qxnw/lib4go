package logger

import (
	"fmt"
	"os"
	"time"
)

func sysLogWrite(level string, content interface{}) {
	e := LogEvent{}
	e.Now = time.Now()
	e.Level = level
	e.Name = "sys"
	e.Session = getSessionID()
	e.Content = fmt.Sprintf("%v", content)
	e.Output = "[%datetime][%l][%session] %content"
	os.Stderr.WriteString(transform(e.Output, e))
}
