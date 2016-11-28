package webserver

import (
	"net/http"
	"time"

	"github.com/arsgo/lib4go/logger"
	"github.com/qxnw/lib4go/utility"
)

//Context 上下文
type Context struct {
	StartTime time.Time
	Writer    http.ResponseWriter
	Request   *http.Request
	Session   string
	Address   string
	Script    string
	Encoding  string
	Log       logger.ILogger
}

//NewContext 构建web请求上下文
func NewContext(loggerName string, w http.ResponseWriter, r *http.Request, address string, script string) *Context {
	context := &Context{Writer: w, Request: r, Address: address, Script: script}
	context.StartTime = time.Now()
	context.Session = utility.GetSessionID()
	context.Log, _ = logger.NewSession(loggerName, context.Session)
	return context

}

//PassTime 计算当前使用已过去的时间
func (c *Context) PassTime() time.Duration {
	return time.Now().Sub(c.StartTime)
}
