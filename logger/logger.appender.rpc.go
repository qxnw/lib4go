package logger

import (
	"bytes"
	"io"
	"sync"
	"time"
)

//RPCAppender 文件输出器
type RPCAppender struct {
	name      string
	buffer    *bytes.Buffer
	lastWrite time.Time
	layout    *Appender
	ticker    *time.Ticker
	locker    sync.Mutex
	writer    io.WriteCloser
	Level     int
}

//NewRPCAppender 构建基于文件流的日志输出对象
func NewRPCAppender(layout *Appender) (fa *RPCAppender, err error) {
	fa = &RPCAppender{layout: layout}
	fa.Level = getLevel(layout.Level)
	fa.ticker = time.NewTicker(TimeWriteToFile)
	fa.writer = nil //构建rpc请求
	go fa.writeTo()
	return
}

//Write 写入日志
func (f *RPCAppender) Write(event LogEvent) {
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
func (f *RPCAppender) Close() {
	f.Level = ILevel_OFF
	f.ticker.Stop()
	f.locker.Lock()
	f.buffer.WriteTo(f.writer)
	f.writer.Close()
	f.locker.Unlock()
}

//writeTo 定时写入文件
func (f *RPCAppender) writeTo() {
START:
	for {
		select {
		case _, ok := <-f.ticker.C:
			if ok {
				f.locker.Lock()
				f.buffer.WriteTo(f.writer)
				f.locker.Unlock()
			} else {
				break START
			}
		}
	}
}
