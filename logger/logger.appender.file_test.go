package logger

import (
	"testing"
	"time"
)

func TestNewFailAppender(t *testing.T) {
	// 创建日志输出对象，创建文件成功
	path := "../logs/test.log"
	layout := &Appender{Type: "file", Level: "All"}
	_, err := NewFileAppender(path, layout)
	if err != nil {
		t.Errorf("test fail:%v", err)
	}

	// 创建日志输出对象，创建文件失败
	path = "/root/test.log"
	layout = &Appender{Type: "file", Level: "All"}
	_, err = NewFileAppender(path, layout)
	if err == nil {
		t.Error("test fail")
	}
}

func TestWrite(t *testing.T) {
	path := "../logs/test.log"
	layout := &Appender{Type: "file", Level: "All"}
	f, err := NewFileAppender(path, layout)
	if err != nil {
		t.Errorf("test fail:%v", err)
	}

	event := LogEvent{Level: "All", Output: "output"}
	f.Write(event)

	// 不能写日志
	f.Level = getLevel("Off")
	f.Write(event)
}

func TestWriteToFileAndReadCheck(tx *testing.T) {
	// 写入文件中
	t, err := time.Parse("2006/01/02 15:04:05", "2016/11/28 16:38:27")
	if err != nil {
		tx.Errorf("test fail, %+v", err)
	}
	path := "../log/20161128.log"
	layout := &Appender{Type: "file", Level: "debug", Path: path}
	fa, err := NewFileAppender(path, layout)
	if err != nil {
		tx.Errorf("test fail:%v", err)
	}
	event := []LogEvent{
		LogEvent{Level: "debug", Now: t, Name: "test", Session: "12345678", Content: "content", Output: "output"},
	}

	// 读取文件，进行校验
}
