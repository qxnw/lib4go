package logger

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"errors"

	"time"

	"github.com/qxnw/lib4go/file"
)

//Appender 输出器
type Appender struct {
	Type   string `json:"type"`
	Level  string `json:"level"`
	Path   string `json:"path"`
	Layout string `json:"layout"`
	Server string `json:"server"`
	Method string `json:"method"`
	Flush  int    `json:"flush"`
}

var loggerPath = file.GetAbs("../conf/ars.logger.json")

//ReadConfig 读取配置文件
func ReadConfig() (appenders []*Appender) {
	var err error
	appenders, err = read()
	if err == nil {
		return
	}
	appenders = getDefConfig()
	sysLoggerError(err)
	err = writeToFile(loggerPath, appenders)
	if err != nil {
		sysLoggerError(err)
	}

	return
}

// // TimeClear 定时清理loggermanager时间间隔
// var TimeClear = time.Second

// TimeWriteToFile 定时写入文件时间间隔
var TimeWriteToFile = time.Second

func read() (appenders []*Appender, err error) {
	appenders = make([]*Appender, 0, 2)
	if !exists(loggerPath) {
		err = errors.New("配置文件不存在:" + loggerPath)
		return
	}
	bytes, err := ioutil.ReadFile(loggerPath)
	if err != nil {
		err = errors.New("无法读取配置文件")
		return
	}
	if err = json.Unmarshal(bytes, &appenders); err != nil {
		err = errors.New("配置文件格式有误，无法序列化")
		return
	}
	return
}
func writeToFile(loggerPath string, appenders []*Appender) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	fwriter, err := file.CreateFile(loggerPath)
	if err != nil {
		return
	}
	data, err := json.Marshal(appenders)
	if err != nil {
		return
	}
	_, err = fwriter.Write(data)
	if err != nil {
		return
	}
	fwriter.Close()
	sysLoggerError("已创建日志配置文件:", loggerPath)
	return
}
func getDefConfig() (appenders []*Appender) {
	fileAppender := &Appender{Type: "file", Level: SLevel_ALL}
	fileAppender.Path = file.GetAbs("../logs/%name/%level_%date.log")
	fileAppender.Layout = "[%datetime][%l][%session] %content%n"
	appenders = append(appenders, fileAppender)

	sdtoutAppender := &Appender{Type: "stdout", Level: SLevel_ALL}
	sdtoutAppender.Layout = "[%datetime][%l][%session] %content%n"
	appenders = append(appenders, sdtoutAppender)

	return
}
func exists(p string) bool {
	_, err := os.Stat(p)
	return err == nil || os.IsExist(err)
}
func getCaller(index int) string {
	defer recover()
	_, file, line, _ := runtime.Caller(index)
	return fmt.Sprintf("%s%d", filepath.Base(file), line)
}
