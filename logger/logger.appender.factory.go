package logger

import (
	"fmt"
	"strings"
)

//makeAppender 根据appender配置及日志信息生成appender对象
func makeAppender(l *Appender, event LogEvent) (IAppender, error) {
	switch strings.ToLower(l.Type) {
	case appender_file:
		return NewFileAppender(makeUniq(l, event), l)
	case appender_stdout:
		return NewStudoutAppender(makeUniq(l, event), l)
	}
	return nil, fmt.Errorf("不支持的日志类型:%s", l.Type)
}

//makeUniq 根据appender配置及日志信息生成appender唯一编号
func makeUniq(l *Appender, event LogEvent) string {
	switch strings.ToLower(l.Type) {
	case appender_file:
		return transform(l.Path, event)
	}
	return l.Type
}
