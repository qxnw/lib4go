/*
	关于Http Server的文档：
		http://www.cnblogs.com/yjf512/archive/2012/08/22/2650873.html
*/

package webserver

import (
	"net"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/qxnw/lib4go/logger"
	"github.com/qxnw/lib4go/utility"
)

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

func NewContext(loggerName string, w http.ResponseWriter, r *http.Request, address string, script string) *Context {
	context := &Context{Writer: w, Request: r, Address: address, Script: script}
	context.StartTime = time.Now()
	context.Session = utility.GetSessionID()
	context.Log, _ = logger.NewSession(loggerName, context.Session)
	return context

}

func (c *Context) PassTime() time.Duration {
	return time.Now().Sub(c.StartTime)
}

//WebHandler Web处理程序
type WebHandler struct {
	LoggerName string
	Path       string
	Script     string
	Method     string
	Encoding   string
	Handler    func(*Context)
}

//WebServer WEB服务
type WebServer struct {
	routes     []WebHandler
	address    string
	loggerName string
	Log        logger.ILogger
	l          net.Listener
}

//NewWebServer 创建WebServer服务
func NewWebServer(address string, loggerName string, handlers ...WebHandler) (server *WebServer) {
	server = &WebServer{routes: handlers, address: address, loggerName: loggerName}
	server.Log, _ = logger.Get(loggerName)
	return
}
func (w WebHandler) recover(log logger.ILogger) {
	if r := recover(); r != nil {
		log.Fatal(r, string(debug.Stack()))
	}
}
func (h WebHandler) call(w http.ResponseWriter, r *http.Request) {
	context := NewContext(h.LoggerName, w, r, h.Path, h.Script)
	defer h.recover(context.Log)
	context.Encoding = h.Encoding
	if strings.EqualFold(h.Method, "*") || strings.EqualFold(r.Method, h.Method) {
		h.Handler(context)
		return
	}
	w.WriteHeader(404)
	w.Write([]byte("您访问的页面不存在"))
}

//Serve 启动WEB服务器
func (w *WebServer) Serve() (err error) {
	mux := http.NewServeMux()
	for _, handler := range w.routes {
		mux.HandleFunc(handler.Path, handler.call)

	}
	l, err := net.Listen("tcp", w.address)
	if err != nil {
		return
	}
	w.l = l
	err = http.Serve(w.l, mux)
	return
}

//Stop 停止服务器
func (w *WebServer) Stop() {
	if w.l != nil {
		w.l.Close()
	}
}
