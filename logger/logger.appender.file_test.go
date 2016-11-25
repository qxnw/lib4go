package logger

import (
	"bytes"
	"io"
	"sync"
	"time"
)

type TestAppender struct {
	name      string
	buffer    *bytes.Buffer
	lastWrite time.Time
	layout    *Appender
	file      io.WriteCloser
	ticker    *time.Ticker
	locker    sync.Mutex
	Level     int
}

func NewTestAppender(path string, layout *Appender) (fa *TestAppender, err error) {
	fa = &TestAppender{layout: layout}
	fa.Level = getLevel(layout.Layout)
	fa.buffer = bytes.NewBufferString("begin")
	fa.ticker = time.NewTicker(time.Second)

	return
}

func (f *TestAppender) Write(event LogEvent) {
	current := getLevel(event.Level)
	if current < f.Level {
		return
	}
	f.locker.Lock()
	f.buffer.WriteString(event.Output)
	f.locker.Unlock()
	f.lastWrite = time.Now()
}

func (f *TestAppender) Close() {
	f.Level = ILevel_OFF
	f.ticker.Stop()
	f.buffer.WriteString("end")
	f.buffer.Reset()
	f.locker.Unlock()
}

func (f *TestAppender) writeTo() {
	for {
		select {
		case _, ok := <-f.ticker.C:
			if ok {
				// f.locker.Lock()
			}
		}
	}
}
