package logger

import (
	"github.com/qxnw/lib4go/concurrent"
)

type loggerManager struct {
	appenders cmap.ConcurrentMap
	configs   []*Appender
}

func newLoggerManager() (m *loggerManager) {
	m = &loggerManager{}
	m.appenders = cmap.New()
	m.configs = ReadConfig()
	return m
}

//Log 将日志内容写入appender, 如果appender不存在则创建
func (a *loggerManager) Log(event LogEvent) {
	go func() {
		for _, config := range a.configs {
			uniq := makeUniq(config, event)
			_, currentAppender, err := a.appenders.SetIfAbsentCb(uniq, func(p ...interface{}) (interface{}, error) {
				l := p[0].(*Appender)
				return makeAppender(l, event)
			}, config)
			if err == nil {
				capp := currentAppender.(IAppender)
				event.Output = transform(event.Content, event)
				capp.Write(event)
			} else {
				sysLogWrite(SLevel_Error, err.Error())
			}
		}
	}()
}
