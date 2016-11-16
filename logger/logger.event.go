package logger

import "time"

type LogEvent struct {
	Level   string
	Now     time.Time
	Name    string
	Session string
	Content string
	Output  string
}

func NewLogEvent(name string, level string, session string, content string) LogEvent {
	e := LogEvent{}
	e.Now = time.Now()
	e.Level = level
	e.Name = name
	e.Session = session
	e.Content = content
	return e

}
