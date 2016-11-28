package logger

import (
	"testing"
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
