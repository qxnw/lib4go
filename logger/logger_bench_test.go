package logger

import (
	"testing"
	"time"
)

// BenchmarkManagerLog 写日志到文件的性能测试
// 目前测试结果：BenchmarkManagerLog-4	 	20000    67800 ns/op
func BenchmarkManagerLog(b *testing.B) {
	manager = newLoggerManager()
	event := LogEvent{Level: "Debug", Now: time.Now(), Name: "benchmark", Session: "12345678", Content: "content1", Output: "output1"}
	for i := 0; i < b.N; i++ {
		manager.Log(event)
	}
}
