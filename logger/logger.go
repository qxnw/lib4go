package logger

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"bytes"

	"github.com/qxnw/lib4go/concurrent/cmap"
	"github.com/qxnw/lib4go/net"
	"github.com/qxnw/lib4go/utility"
)

//Logger 日志对象
type Logger struct {
	index    int64
	names    string
	sessions string
	tags     map[string]string
}
type event struct {
	f       int
	tp      string
	fm      string
	name    string
	session string
	value   []interface{}
}

var loggerEventChan chan *LogEvent
var loggerCloserChan chan *Logger
var loggerPool *sync.Pool
var loggers cmap.ConcurrentMap
var manager *loggerManager
var LocalIP string

func init() {
	loggerPool = &sync.Pool{
		New: func() interface{} {
			return New("")
		},
	}
	LocalIP = net.GetLocalIPAddress()
	register(appender_file, readFromFile)
	var err error
	manager, err = newLoggerManager()
	if err != nil {
		fmt.Println("logger err:未启用日志")
		return
	}
	loggerEventChan = make(chan *LogEvent, 2000)
	loggerCloserChan = make(chan *Logger, 1000)
	for i := 0; i < 50; i++ {
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
	logger = &Logger{index: 100}
	logger.names = names
	logger.sessions = CreateSession()
	return logger
}

//GetSession 根据日志名称及session获取日志组件
func GetSession(name string, sessionID string, tags ...string) (logger *Logger) {
	logger = loggerPool.Get().(*Logger)
	logger.names = name
	logger.sessions = sessionID
	logger.tags = make(map[string]string)
	if len(tags) > 0 && len(tags) != 2 {
		panic(fmt.Sprintf("日志输入参数错误，扩展参数必须成对出现:%s,%v", name, tags))
	}
	for i := 0; i < len(tags)-1; i++ {
		logger.tags[tags[i]] = tags[i+1]
	}
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

//WaitClose 等待所有日志写入完毕并关闭
func (logger *Logger) WaitClose() {
	logger.Close()
	time.Sleep(time.Second * 2)
}

//SetTag 设置tag
func (logger *Logger) SetTag(name string, value string) {
	logger.tags[name] = value
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
	logger.log(SLevel_Debug, content...)
}

//Debugf 输出debug日志
func (logger *Logger) Debugf(format string, content ...interface{}) {
	if !isOpen {
		return
	}
	logger.logfmt(format, SLevel_Debug, content...)
}

//Info 输出info日志
func (logger *Logger) Info(content ...interface{}) {
	if !isOpen {
		return
	}
	logger.log(SLevel_Info, content...)
}

//Infof 输出info日志
func (logger *Logger) Infof(format string, content ...interface{}) {
	if !isOpen {
		return
	}
	logger.logfmt(format, SLevel_Info, content...)
}

//Warn 输出info日志
func (logger *Logger) Warn(content ...interface{}) {
	if !isOpen {
		return
	}
	logger.log(SLevel_Warn, content...)
}

//Warnf 输出info日志
func (logger *Logger) Warnf(format string, content ...interface{}) {
	if !isOpen {
		return
	}
	logger.logfmt(format, SLevel_Warn, content...)
}

//Error 输出Error日志
func (logger *Logger) Error(content ...interface{}) {
	if !isOpen {
		return
	}
	logger.log(SLevel_Error, content...)
}

//Errorf 输出Errorf日志
func (logger *Logger) Errorf(format string, content ...interface{}) {
	if !isOpen {
		return
	}
	logger.logfmt(format, SLevel_Error, content...)
}

//Fatal 输出Fatal日志
func (logger *Logger) Fatal(content ...interface{}) {
	if !isOpen {
		return
	}
	logger.log(SLevel_Fatal, content...)
}

//Fatalf 输出Fatalf日志
func (logger *Logger) Fatalf(format string, content ...interface{}) {
	if !isOpen {
		return
	}
	logger.logfmt(format, SLevel_Fatal, content...)

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
func (logger *Logger) logfmt(f string, level string, content ...interface{}) {
	event := NewLogEvent(logger.names, level, logger.sessions, fmt.Sprintf(f, content...), nil, atomic.AddInt64(&logger.index, 1))
	loggerEventChan <- event
}
func (logger *Logger) log(level string, content ...interface{}) {
	event := NewLogEvent(logger.names, level, logger.sessions, getString(content...), nil, atomic.AddInt64(&logger.index, 1))
	loggerEventChan <- event
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
			v.Close()
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
