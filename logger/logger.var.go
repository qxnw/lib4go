package logger

const (
	ILevel_ALL = iota
	ILevel_Debug
	ILevel_Info
	ILevel_Error
	ILevel_Fatal
	ILevel_OFF
)
const (
	SLevel_OFF   = "Off"
	SLevel_Info  = "Info"
	SLevel_Error = "Error"
	SLevel_Fatal = "Fatal"
	SLevel_Debug = "Debug"
	SLevel_ALL   = "All"
)

const (
	appender_file = "file"
)

var levelMap map[string]int

func init() {
	levelMap = map[string]int{
		SLevel_OFF:   ILevel_OFF,
		SLevel_Info:  ILevel_Info,
		SLevel_Error: ILevel_Error,
		SLevel_Fatal: ILevel_Fatal,
		SLevel_Debug: ILevel_Debug,
		SLevel_ALL:   ILevel_ALL,
	}
}

func getLevel(name string) int {
	if l, ok := levelMap[name]; ok {
		return l
	}
	return ILevel_ALL
}