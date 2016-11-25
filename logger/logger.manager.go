package logger

import "github.com/qxnw/lib4go/concurrent"
import "time"

type loggerManager struct {
	appenders cmap.ConcurrentMap
	factory   *loggerAppenderFactory
	configs   []*Appender
	ticker    *time.Ticker
	isClose   bool
}
type appenderEntity struct {
	appender IAppender
	last     time.Time
}

func newLoggerManager() (m *loggerManager) {
	m = &loggerManager{isClose: false}
	m.factory = &loggerAppenderFactory{}
	m.appenders = cmap.New()
	m.configs = ReadConfig()
	m.ticker = time.NewTicker(time.Second)
	go m.clearUp()
	return m
}

//Log 将日志内容写入appender, 如果appender不存在则创建
func (a *loggerManager) Log(event LogEvent) {
	if a.isClose {
		return
	}
	for _, config := range a.configs {
		uniq := a.factory.MakeUniq(config, event)
		_, currentAppender, err := a.appenders.SetIfAbsentCb(uniq, func(p ...interface{}) (entity interface{}, err error) {
			l := p[0].(*Appender)
			app, err := a.factory.MakeAppender(l, event)
			if err != nil {
				return
			}
			entity = &appenderEntity{appender: app, last: time.Now()}
			return
		}, config)
		if err == nil {
			capp := currentAppender.(*appenderEntity)
			event.Output = transform(config.Layout, event)
			capp.appender.Write(event)
			capp.last = time.Now()
		} else {
			sysLoggerError(err)
		}
	}
}
func (a *loggerManager) clearUp() {
START:
	for {
		select {
		case _, ok := <-a.ticker.C:
			if ok {
				count := a.appenders.RemoveIterCb(func(key string, v interface{}) bool {
					apd := v.(*appenderEntity)
					if time.Now().Sub(apd.last).Seconds() > 10 {
						apd.appender.Close()
						return true
					}
					return false
				})
				if count > 0 {
					sysLoggerInfo("已移除:", count)
				}
			} else {
				break START
			}
		}
	}
}

func (a *loggerManager) Close() {
	a.isClose = true
	a.ticker.Stop()
	a.appenders.RemoveIterCb(func(key string, v interface{}) bool {
		apd := v.(*appenderEntity)
		apd.appender.Close()
		return true
	})
}
