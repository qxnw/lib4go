package logger

import (
	"os"
	"time"

	"sync"

	"github.com/lunny/log"
)

//StdoutAppender 标准输出器
type StdoutAppender struct {
	name      string
	lastWrite time.Time
	layout    *Appender
	output    *log.Logger
	unq       string
	Level     int
	mu        sync.Mutex
}

//NewStudoutAppender 构建基于文件流的日志输出对象
func NewStudoutAppender(unq string, layout *Appender) (fa *StdoutAppender, err error) {
	fa = &StdoutAppender{layout: layout, unq: unq}
	fa.Level = getLevel(layout.Level)
	fa.output = log.New(os.Stdout, "", log.Ldefault())
	fa.output.SetOutputLevel(log.Ldebug)
	return
}

//Write 写入日志
func (f *StdoutAppender) Write(event LogEvent) {
	current := getLevel(event.Level)
	if current < f.Level {
		return
	}
	f.lastWrite = time.Now()
	f.mu.Lock()
	switch current {
	case ILevel_Debug:
		f.output.Debug(event.Output)
	case ILevel_Error:
		f.output.Error(event.Output)
	case ILevel_Info:
		f.output.Info(event.Output)
	case ILevel_Fatal:
		f.output.Fatal(event.Output)
	}
	f.mu.Unlock()
}

//Close 关闭当前appender
func (f *StdoutAppender) Close() {
	f.Level = ILevel_OFF
}
