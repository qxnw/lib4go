package logger

import (
	"bytes"
	"io"
	"sync"
	"time"

	"github.com/qxnw/lib4go/file"
)

//FileAppender 文件输出器
type FileAppender struct {
	name      string
	buffer    *bytes.Buffer
	lastWrite time.Time
	layout    *Appender
	file      io.WriteCloser
	ticker    *time.Ticker
	locker    sync.Mutex
	Level     int
}

//NewFileAppender 构建基于文件流的日志输出对象
func NewFileAppender(path string, layout *Appender) (fa *FileAppender, err error) {
	fa = &FileAppender{layout: layout}
	fa.Level = getLevel(layout.Level)
	fa.buffer = bytes.NewBufferString("\n---------------------begin-------------------------\n\n")
	fa.ticker = time.NewTicker(time.Second)
	fa.file, err = file.CreateFile(path)
	if err != nil {
		return
	}
	go fa.writeTo()
	return
}

//Write 写入日志
func (f *FileAppender) Write(event LogEvent) {
	current := getLevel(event.Level)
	if current < f.Level {
		return
	}
	f.locker.Lock()
	f.buffer.WriteString(event.Output)
	f.locker.Unlock()
	f.lastWrite = time.Now()
}

//Close 关闭当前appender
func (f *FileAppender) Close() {
	f.Level = ILevel_OFF
	f.ticker.Stop()
	f.locker.Lock()
	f.buffer.WriteString("\n---------------------end-------------------------\n")
	f.buffer.WriteTo(f.file)
	f.file.Close()
	f.locker.Unlock()
}

//writeTo 定时写入文件
func (f *FileAppender) writeTo() {
START:
	for {
		select {
		case _, ok := <-f.ticker.C:
			if ok {
				f.locker.Lock()
				f.buffer.WriteTo(f.file)
				f.locker.Unlock()
			} else {
				break START
			}
		}
	}
}
