package logger

import (
	"fmt"
	"sync"

	"bytes"

	"github.com/qxnw/lib4go/concurrent/cmap"
	"github.com/qxnw/lib4go/utility"
)

//Logger 日志对象
type Logger struct {
	names    string
	sessions string
}
type event struct {
	f       int
	tp      string
	fm      string
	name    string
	session string
	value   []interface{}
}

var loggerEventChan chan event
var loggerCloserChan chan *Logger
var loggerPool *sync.Pool
var loggers cmap.ConcurrentMap
var manager *loggerManager

func init() {
	loggerPool = &sync.Pool{
		New: func() interface{} {
			return New("")
		},
	}

	register(appender_file, readFromFile)
	var err error
	manager, err = newLoggerManager()
	if err != nil {
		fmt.Println("logger err:未启用日志")
		return
	}
	loggerEventChan = make(chan event, 2000)
	loggerCloserChan = make(chan *Logger, 1000)
	for i := 0; i < 100; i++ {
		go logNow()
	}
}

//ResetConfig 重置日志配置
func ResetConfig(conf string) (err error) {
	apds, err := NewAppender(conf)
	if err != nil {
		return err
	}
	manager.configs = apds
	return nil
}

//New 根据一个或多个日志名称构建日志对象，该日志对象具有新的session id系统不会缓存该日志组件
func New(names string) (logger *Logger) {
	logger = &Logger{}
	logger.names = names
	logger.sessions = CreateSession()
	return logger
}

//GetSession 根据日志名称及session获取日志组件
func GetSession(name string, sessionID string) (logger *Logger) {
	logger = loggerPool.Get().(*Logger)
	logger.names = name
	logger.sessions = sessionID
	return logger
}

//Close 关闭当前日志组件
func (logger *Logger) Close() {
	select {
	case loggerCloserChan <- logger:
	default:
		loggerPool.Put(logger)
	}
}

//SetTag 设置tag
func (logger *Logger) SetTag(name string, value string) {
	logger.SetTag(name, value)
}

//GetSessionID 获取当前日志的session id
func (logger *Logger) GetSessionID() string {
	if len(logger.sessions) > 0 {
		return logger.sessions
	}
	return ""
}

//Debug 输出debug日志
func (logger *Logger) Debug(content ...interface{}) {
	if !isOpen {
		return
	}
	loggerEventChan <- event{f: 1, tp: SLevel_Debug, name: logger.names, session: logger.sessions, value: content}
}

//Debugf 输出debug日志
func (logger *Logger) Debugf(format string, content ...interface{}) {
	if !isOpen {
		return
	}
	loggerEventChan <- event{f: 0, fm: format, tp: SLevel_Debug, name: logger.names, session: logger.sessions, value: content}
}

//Info 输出info日志
func (logger *Logger) Info(content ...interface{}) {
	if !isOpen {
		return
	}
	loggerEventChan <- event{f: 1, tp: SLevel_Info, name: logger.names, session: logger.sessions, value: content}
}

//Infof 输出info日志
func (logger *Logger) Infof(format string, content ...interface{}) {
	if !isOpen {
		return
	}
	loggerEventChan <- event{f: 0, fm: format, tp: SLevel_Info, name: logger.names, session: logger.sessions, value: content}
}

//Warn 输出info日志
func (logger *Logger) Warn(content ...interface{}) {
	if !isOpen {
		return
	}
	loggerEventChan <- event{f: 1, tp: SLevel_Warn, name: logger.names, session: logger.sessions, value: content}
}

//Warnf 输出info日志
func (logger *Logger) Warnf(format string, content ...interface{}) {
	if !isOpen {
		return
	}
	loggerEventChan <- event{f: 0, fm: format, tp: SLevel_Warn, name: logger.names, session: logger.sessions, value: content}
}

//Error 输出Error日志
func (logger *Logger) Error(content ...interface{}) {
	if !isOpen {
		return
	}
	loggerEventChan <- event{f: 1, tp: SLevel_Error, name: logger.names, session: logger.sessions, value: content}

}

//Errorf 输出Errorf日志
func (logger *Logger) Errorf(format string, content ...interface{}) {
	if !isOpen {
		return
	}
	loggerEventChan <- event{f: 0, fm: format, tp: SLevel_Error, name: logger.names, session: logger.sessions, value: content}
}

//Fatal 输出Fatal日志
func (logger *Logger) Fatal(content ...interface{}) {
	if !isOpen {
		return
	}
	loggerEventChan <- event{f: 1, tp: SLevel_Fatal, name: logger.names, session: logger.sessions, value: content}
}

//Fatalf 输出Fatalf日志
func (logger *Logger) Fatalf(format string, content ...interface{}) {
	if !isOpen {
		return
	}
	loggerEventChan <- event{f: 0, fm: format, tp: SLevel_Fatal, name: logger.names, session: logger.sessions, value: content}

}

//Fatalln 输出Fatal日志
func (logger *Logger) Fatalln(content ...interface{}) {
	logger.Fatal(content...)
}

//Print 输出info日志
func (logger *Logger) Print(content ...interface{}) {
	logger.Info(content...)

}

//Printf 输出info日志
func (logger *Logger) Printf(format string, content ...interface{}) {
	if logger == nil {
		return
	}
	logger.Infof(format, content...)
}

//Println 输出info日志
func (logger *Logger) Println(content ...interface{}) {
	logger.Print(content...)

}
func logNow() {
	for {
		select {
		case logger := <-loggerCloserChan:
			loggerPool.Put(logger)
		case v, ok := <-loggerEventChan:
			if !ok {
				return
			}
			if v.f == 1 {
				event := NewLogEvent(v.name, v.tp, v.session, getString(v.value...), nil)
				manager.Log(event)
				continue
			}
			event := NewLogEvent(v.name, v.tp, v.session, fmt.Sprintf(v.fm, v.value...), nil)
			manager.Log(event)
		}
	}
}
func getString(c ...interface{}) string {
	if len(c) == 1 {
		return fmt.Sprintf("%v", c[0])
	}
	var buf bytes.Buffer
	for i := 0; i < len(c); i++ {
		buf.WriteString(fmt.Sprint(c[i]))
		if i != len(c)-1 {
			buf.WriteString(" ")
		}
	}
	return buf.String()
}

//Close 关闭所有日志组件
func Close() {
	if manager == nil {
		return
	}
	manager.Close()
	close(loggerEventChan)
}

//CreateSession create logger session
func CreateSession() string {
	return utility.GetGUID()[0:9]
}
