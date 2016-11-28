package logger

import (
	"testing"
	"time"
)

func TestLog(t *testing.T) {
	manager.isClose = true
	manager.Log(LogEvent{Level: "Info", Now: time.Now(), Name: "name1", Session: "12345678", Content: "content1", Output: "output1"})

	// 测试完成打开appender，否则音响其他测试
	manager.isClose = false
}

func TestManagerClose(t *testing.T) {
	manager.Log(LogEvent{Level: "Info", Now: time.Now(), Name: "name1", Session: "12345678", Content: "content1", Output: "output1"})
	if len(manager.appenders.Keys()) == 0 {
		t.Error("test fail")
	}

	manager.Close()
	if len(manager.appenders.Keys()) != 0 {
		t.Errorf("test fail:manager appenders have:%d", len(manager.appenders.Keys()))
	}
}
