package logger

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
	return configAdapter[defaultConfigAdapter]()
}
