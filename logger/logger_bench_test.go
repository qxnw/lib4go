package logger

import (
	"testing"
	"time"
)

func BenchmarkManagerLog(b *testing.B) {
	manager = newLoggerManager()
	event := LogEvent{Level: "Debug", Now: time.Now(), Name: "benchmark", Session: "12345678", Content: "content1", Output: "output1"}
	for i := 0; i < b.N; i++ {
		manager.Log(event)
	}
}
