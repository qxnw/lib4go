package logger

import (
	"fmt"
	"os"

	"github.com/qxnw/lib4go/concurrent/cmap"
	"github.com/qxnw/lib4go/utility"
)

//Logger 日志对象
type Logger struct {
	names    []string
	sessions []string
	tags     map[string]string
}

var loggers cmap.ConcurrentMap
var manager *loggerManager

func init() {
	register(appender_file, readFromFile)
	loggers = cmap.New()
	manager = newLoggerManager()
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
	logger.tags = make(map[string]string)
	for range names {
		logger.sessions = append(logger.sessions, getSessionID())
	}
	return logger
}

//Get 根据名称从缓存中获取日志组件，如果缓存中不存在则创建并保存到缓存中
func Get(names ...string) (logger *Logger) {
	logger = &Logger{}
	for _, name := range names {
		_, session := loggers.SetIfAbsent(name, getSessionID())
		logger.names = append(logger.names, name)
		logger.sessions = append(logger.sessions, session.(string))
	}
	return logger
}

//GetSession 根据日志名称及session获取日志组件
func GetSession(name string, sessionID string) (logger *Logger) {
	logger = &Logger{}
	logger.names = append(logger.names, name)
	logger.sessions = append(logger.sessions, sessionID)
	return logger
}

//SetTag 设置tag
func (logger *Logger) SetTag(name string, value string) {
	logger.tags[name] = value
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
	for i, name := range logger.names {
		event := NewLogEvent(name, SLevel_Debug, logger.sessions[i], fmt.Sprint(content...), logger.tags)
		go manager.Log(event)
	}
}

//Debugf 输出debug日志
func (logger *Logger) Debugf(format string, content ...interface{}) {
	for i, name := range logger.names {
		event := NewLogEvent(name, SLevel_Debug, logger.sessions[i], fmt.Sprintf(format, content...), logger.tags)
		go manager.Log(event)
	}
}

//Info 输出info日志
func (logger *Logger) Info(content ...interface{}) {
	for i, name := range logger.names {
		event := NewLogEvent(name, SLevel_Info, logger.sessions[i], fmt.Sprint(content...), logger.tags)
		go manager.Log(event)
	}
}

//Infof 输出info日志
func (logger *Logger) Infof(format string, content ...interface{}) {
	for i, name := range logger.names {
		event := NewLogEvent(name, SLevel_Info, logger.sessions[i], fmt.Sprintf(format, content...), logger.tags)
		go manager.Log(event)
	}
}

//Error 输出Error日志
func (logger *Logger) Error(content ...interface{}) {
	for i, name := range logger.names {
		event := NewLogEvent(name, SLevel_Error, logger.sessions[i], fmt.Sprint(content...), logger.tags)
		go manager.Log(event)
	}

}

//Errorf 输出Errorf日志
func (logger *Logger) Errorf(format string, content ...interface{}) {
	for i, name := range logger.names {
		event := NewLogEvent(name, SLevel_Error, logger.sessions[i], fmt.Sprintf(format, content...), logger.tags)
		go manager.Log(event)
	}
}

//Fatal 输出Fatal日志
func (logger *Logger) Fatal(content ...interface{}) {
	for i, name := range logger.names {
		event := NewLogEvent(name, SLevel_Fatal, logger.sessions[i], fmt.Sprint(content...), logger.tags)
		go manager.Log(event)
	}
	os.Exit(999)

}

//Fatalf 输出Fatalf日志
func (logger *Logger) Fatalf(format string, content ...interface{}) {
	for i, name := range logger.names {
		event := NewLogEvent(name, SLevel_Fatal, logger.sessions[i], fmt.Sprintf(format, content...), logger.tags)
		go manager.Log(event)
	}
	os.Exit(999)
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
	logger.Infof(format, content...)
}

//Println 输出info日志
func (logger *Logger) Println(content ...interface{}) {
	logger.Print(content...)

}
func getSessionID() string {
	id := utility.GetGUID()
	return id[:8]
}

//Close 关闭所有日志组件
func Close() {
	manager.Close()
}
