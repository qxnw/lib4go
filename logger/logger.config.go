package logger

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"errors"

	"github.com/qxnw/lib4go/utility"
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

//ReadConfig 读取配置文件
func ReadConfig() (appenders []*Appender) {
	var err error
	appenders, err = read()
	if err != nil {
		appenders = getDefConfig()
		sysLogWrite(SLevel_Error, err)
	}
	return
}

func read() (appenders []*Appender, err error) {
	loggerPath := utility.GetExcPath("./conf/ars.logger.json", "bin")
	appenders = make([]*Appender, 0, 2)
	if !exists(loggerPath) {
		err = errors.New("配置文件不存在")
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
func writeToFile(loggerPath string, appenders []*Appender) {
	if r := recover(); r != nil {
		sysLogWrite(SLevel_Error, r)
	}
	data, _ := json.Marshal(appenders)
	err := ioutil.WriteFile(loggerPath, data, os.ModeAppend)
	if err != nil {
		sysLogWrite(SLevel_Error, err)
	}
	return
}
func getDefConfig() (appenders []*Appender) {
	appender := &Appender{Type: "file", Level: SLevel_ALL}
	appender.Path = utility.GetExcPath("./logs/%name/%level_%date.log", "bin")
	appender.Layout = "[%datetime][%l][%session] %content"
	appenders = append(appenders, appender)
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
