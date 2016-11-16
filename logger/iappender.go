package logger

type IAppender interface {
	Write(LogEvent)
	Close()
}
