package logger

import "time"
import "sync"
import "github.com/qxnw/lib4go/concurrent/cmap"

type LogEvent struct {
	Level   string
	Now     time.Time
	Name    string
	Session string
	Content string
	Output  string
	Tags    cmap.ConcurrentMap
	lk      sync.Mutex
}

func NewLogEvent(name string, level string, session string, content string, tags cmap.ConcurrentMap) LogEvent {
	e := LogEvent{}
	e.Now = time.Now()
	e.Level = level
	e.Name = name
	e.Session = session
	e.Content = content
	e.Tags = tags
	e.Tags.Set("caller", getCaller(5))
	return e
}
