package logger

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"github.com/qxnw/lib4go/utility"
)

type LoggerConfig struct {
	Layout    string      `json:"layout`
	Appenders []*Appender `json:"appenders"`
}

type Appender struct {
	Type   string `json:"type"`
	Level  string `json:"level"`
	Path   string `json:"path"`
	Output string `json:"output"`
	Server string `json:"server"`
	Method string `json:"method"`
	Flush  int    `json:"flush"`
}

var loggerPath string = utility.GetExcPath("./conf/ars.logger.json", "bin")

func ReadConfig() (config *LoggerConfig, err error) {
	config = &LoggerConfig{}
	if !exists(loggerPath) {
		ad := &Appender{Type: "file", Level: SLevel_ALL}
		ad.Path = utility.GetExcPath("./logs/%name/%level_%date.log", "bin")
		config.Appenders = make([]*Appender, 0, 1)
		config.Appenders = append(config.Appenders, ad)
		config.Layout = "[%datetime][%l][%session] %content"
		data, _ := json.Marshal(config)
		err = ioutil.WriteFile(loggerPath, data, os.ModeAppend)
		return
	}
	bytes, err := ioutil.ReadFile(loggerPath)
	if err != nil {
		return
	}
	if err = json.Unmarshal(bytes, &config); err != nil {
		return
	}
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
