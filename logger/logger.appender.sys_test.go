package logger

import (
	"testing"
)

func TestSysLoogerInfo(t *testing.T) {
	sysLoggerInfo(nil, "content")
	sysLoggerInfo(nil, nil)
	sysLoggerInfo(nil, []string{"1", "2"})
	sysLoggerInfo(nil, [2]string{"1", "2"})
	sysLoggerInfo(nil, TestType{name: "name", age: 12})
	sysLoggerInfo(nil, "content")
}

func TestLoggerError(t *testing.T) {
	sysLoggerError(nil, "content")
	sysLoggerError(nil, nil)
	sysLoggerError(nil, []string{"1", "2"})
	sysLoggerError(nil, [2]string{"1", "2"})
	sysLoggerError(nil, TestType{name: "name", age: 12})
	sysLoggerError(nil, "content")
}

func TestSysLoggerWrite(t *testing.T) {
	sysLoggerWrite(nil, "info", "content")
	sysLoggerWrite(nil, "test", nil)
	sysLoggerWrite(nil, "info", []string{"1", "2"})
	sysLoggerWrite(nil, "info", [2]string{"1", "2"})
	sysLoggerWrite(nil, "info", TestType{name: "name", age: 12})
	sysLoggerWrite(nil, "", "content")
}
