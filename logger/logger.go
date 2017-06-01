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
	names    []string
	sessions []string
	tags     cmap.ConcurrentMap
}

var loggerEventChan chan LogEvent
var loggerCloserChan chan *Logger
var loggerPool *sync.Pool
var loggers cmap.ConcurrentMap
var manager *loggerManager

func init() {
	loggerPool = &sync.Pool{
		New: func() interface{} {
			return New()
		},
	}

	register(appender_file, readFromFile)
	var err error
	manager, err = newLoggerManager()
	if err != nil {
		fmt.Println("logger err:未启用日志")
		return
	}
	loggerEventChan = make(chan LogEvent, 2000)
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
func New(names ...string) (logger *Logger) {
	logger = &Logger{}
	logger.names = names
	logger.tags = cmap.New()
	for range names {
		logger.sessions = append(logger.sessions, CreateSession())
	}
	return logger
}

//GetSession 根据日志名称及session获取日志组件
func GetSession(name string, sessionID string) (logger *Logger) {
	logger = loggerPool.Get().(*Logger)
	logger.names = []string{name}
	logger.sessions = []string{sessionID}
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
		return logger.sessions[0]
	}
	return ""
}

//Debug 输出debug日志
func (logger *Logger) Debug(content ...interface{}) {
	if !isOpen {
		return
	}
	for i, name := range logger.names {
		event := NewLogEvent(name, SLevel_Debug, logger.sessions[i], getString(content...), logger.tags)
		loggerEventChan <- event
	}
}

//Debugf 输出debug日志
func (logger *Logger) Debugf(format string, content ...interface{}) {
	if !isOpen {
		return
	}
	for i, name := range logger.names {
		event := NewLogEvent(name, SLevel_Debug, logger.sessions[i], fmt.Sprintf(format, content...), logger.tags)
		loggerEventChan <- event
	}
}

//Info 输出info日志
func (logger *Logger) Info(content ...interface{}) {
	if !isOpen {
		return
	}
	for i, name := range logger.names {
		event := NewLogEvent(name, SLevel_Info, logger.sessions[i], getString(content...), logger.tags)
		loggerEventChan <- event
	}
}

//Infof 输出info日志
func (logger *Logger) Infof(format string, content ...interface{}) {
	if !isOpen {
		return
	}
	for i, name := range logger.names {
		event := NewLogEvent(name, SLevel_Info, logger.sessions[i], fmt.Sprintf(format, content...), logger.tags)
		loggerEventChan <- event
	}
}

//Warn 输出info日志
func (logger *Logger) Warn(content ...interface{}) {
	if !isOpen {
		return
	}
	for i, name := range logger.names {
		event := NewLogEvent(name, SLevel_Warn, logger.sessions[i], getString(content...), logger.tags)
		loggerEventChan <- event
	}
}

//Warnf 输出info日志
func (logger *Logger) Warnf(format string, content ...interface{}) {
	if !isOpen {
		return
	}
	for i, name := range logger.names {
		event := NewLogEvent(name, SLevel_Warn, logger.sessions[i], fmt.Sprintf(format, content...), logger.tags)
		loggerEventChan <- event
	}
}

//Error 输出Error日志
func (logger *Logger) Error(content ...interface{}) {
	if !isOpen {
		return
	}
	for i, name := range logger.names {
		event := NewLogEvent(name, SLevel_Error, logger.sessions[i], getString(content...), logger.tags)
		loggerEventChan <- event
	}

}

//Errorf 输出Errorf日志
func (logger *Logger) Errorf(format string, content ...interface{}) {
	if !isOpen {
		return
	}
	for i, name := range logger.names {
		event := NewLogEvent(name, SLevel_Error, logger.sessions[i], fmt.Sprintf(format, content...), logger.tags)
		loggerEventChan <- event
	}
}

//Fatal 输出Fatal日志
func (logger *Logger) Fatal(content ...interface{}) {
	if !isOpen {
		return
	}
	for i, name := range logger.names {
		event := NewLogEvent(name, SLevel_Fatal, logger.sessions[i], getString(content...), logger.tags)
		loggerEventChan <- event
	}
}

//Fatalf 输出Fatalf日志
func (logger *Logger) Fatalf(format string, content ...interface{}) {
	if !isOpen {
		return
	}
	for i, name := range logger.names {
		event := NewLogEvent(name, SLevel_Fatal, logger.sessions[i], fmt.Sprintf(format, content...), logger.tags)
		loggerEventChan <- event
	}
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
			manager.Log(v)
		}
	}
}
func getString(c ...interface{}) string {
	var buf bytes.Buffer
	for i := 0; i < len(c); i++ {
		buf.WriteString(fmt.Sprintf("%v", c[i]))
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
