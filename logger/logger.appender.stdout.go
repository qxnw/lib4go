package logger

import (
	"os"
	"time"
)

//StdoutAppender 标准输出器
type StdoutAppender struct {
	name      string
	lastWrite time.Time
	layout    *Appender
	unq       string
	Level     int
}

//NewStudoutAppender 构建基于文件流的日志输出对象
func NewStudoutAppender(unq string, layout *Appender) (fa *StdoutAppender, err error) {
	fa = &StdoutAppender{layout: layout, unq: unq}
	fa.Level = getLevel(layout.Level)
	return
}

//Write 写入日志
func (f *StdoutAppender) Write(event LogEvent) {
	current := getLevel(event.Level)
	if current < f.Level {
		return
	}
	f.lastWrite = time.Now()
	os.Stdout.WriteString(event.Output)
}

//Close 关闭当前appender
func (f *StdoutAppender) Close() {
	f.Level = ILevel_OFF
}
