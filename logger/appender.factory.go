package logger

import (
	"fmt"
	"strings"
)

func makeAppender(l *Appender, event LogEvent) (IAppender, error) {
	switch strings.ToLower(l.Type) {
	case appender_file:
		return NewFileAppender(makeUniq(l, event), l)
	}
	return nil, fmt.Errorf("不支持的日志类型:%s", l.Type)
}

func makeUniq(l *Appender, event LogEvent) string {
	switch strings.ToLower(l.Type) {
	case appender_file:
		return transform(l.Path, event)
	}
	return l.Type
}
