package logger

import (
	"github.com/qxnw/lib4go/concurrent"
)

type appenderManager struct {
	appenders cmap.ConcurrentMap
	config    *LoggerConfig
}

func newAppenderManager() (m *appenderManager) {
	m = &appenderManager{}
	m.appenders = cmap.New()
	m.config, _ = ReadConfig()
	return m
}

//Log 将日志内容写入appender, 如果appender不存在则创建
func (a *appenderManager) Log(event LogEvent) {
	event.Output = transform(a.config.Layout, event)
	go func() {
		for _, config := range a.config.Appenders {
			uniq := makeUniq(config, event)
			_, currentAppender, err := a.appenders.SetIfAbsentCb(uniq, func(p ...interface{}) (interface{}, error) {
				l := p[0].(*Appender)
				return makeAppender(l, event)
			}, config)
			if err == nil {
				capp := currentAppender.(IAppender)
				capp.Write(event)
			}
		}
	}()

}
