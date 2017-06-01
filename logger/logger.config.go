package logger

//Appender 输出器
type Appender struct {
	Type   string `json:"type"`
	Level  string `json:"level"`
	Path   string `json:"path,omitempty"`
	Layout string `json:"layout"`
	Server string `json:"server,omitempty"`
	Method string `json:"method,omitempty"`
	Flush  int    `json:"flush,omitempty"`
}

//ReadConfig 读取配置文件
func ReadConfig() (appenders []*Appender) {
	return configAdapter[defaultConfigAdapter]()
}
