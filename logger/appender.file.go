package logger

import "bytes"
import "os"
import "time"
import "sync"
import "path/filepath"

//FileAppender 文件输出
type FileAppender struct {
	name      string
	buffer    *bytes.Buffer
	lastWrite time.Time
	layout    *Appender
	path      string
	file      *os.File
	ticker    *time.Ticker
	locker    sync.Mutex
	Level     int
}

//NewFileAppender 构建基于文件流的日志输出对象
func NewFileAppender(path string, layout *Appender) (fa *FileAppender, err error) {
	fa = &FileAppender{layout: layout}
	fa.Level = getLevel(layout.Level)
	fa.buffer = bytes.NewBufferString("---------------------begin-------------------------\n")
	fa.ticker = time.NewTicker(time.Second)
	fa.path, err = filepath.Abs(path)
	if err != nil {
		return
	}
	fa.file, err = os.Create(fa.path)
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
	f.locker.Lock()
	f.buffer.WriteString("\n---------------------end-------------------------")
	f.buffer.WriteTo(f.file)
	f.locker.Unlock()
	f.ticker.Stop()
}

//writeTo 定时写入文件
func (f *FileAppender) writeTo() {
	for {
		select {
		case <-f.ticker.C:
			f.locker.Lock()
			f.buffer.WriteTo(f.file)
			f.locker.Unlock()
		}
	}
}
