package logger

import (
	"bytes"
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
	buffer    *bytes.Buffer
	ticker    *time.Ticker
	unq       string
	Level     int
	locker    sync.Mutex
}

//NewStudoutAppender 构建基于文件流的日志输出对象
func NewStudoutAppender(unq string, layout *Appender) (fa *StdoutAppender, err error) {
	fa = &StdoutAppender{layout: layout, unq: unq}
	fa.Level = getLevel(layout.Level)
	fa.buffer = bytes.NewBufferString("")
	fa.output = log.New(fa.buffer, "", log.Llongcolor)
	fa.ticker = time.NewTicker(TimeWriteToSTD)
	fa.output.SetOutputLevel(log.Ldebug)
	go fa.writeTo()
	return
}

//Write 写入日志
func (f *StdoutAppender) Write(event *LogEvent) {
	current := getLevel(event.Level)
	if current < f.Level {
		return
	}
	f.lastWrite = time.Now()
	f.locker.Lock()
	switch current {
	case ILevel_Debug:
		f.output.Debug(event.Output)
	case ILevel_Info:
		f.output.Info(event.Output)
	case ILevel_Warn:
		f.output.Warn(event.Output)
	case ILevel_Error:
		f.output.Error(event.Output)
	case ILevel_Fatal:
		f.output.Fatal(event.Output)
	}
	f.locker.Unlock()
}

//writeTo 定时写入文件
func (f *StdoutAppender) writeTo() {
START:
	for {
		select {
		case _, ok := <-f.ticker.C:
			if ok {
				f.locker.Lock()
				f.buffer.WriteTo(os.Stdout)
				f.buffer.Reset()
				f.locker.Unlock()
			} else {
				break START
			}
		}
	}
}

//Close 关闭当前appender
func (f *StdoutAppender) Close() {
	f.locker.Lock()
	f.Level = ILevel_OFF
	f.ticker.Stop()
	f.locker.Unlock()
}
