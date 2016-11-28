package http

import (
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/arsgo/lib4go/logger"
)

//WebHandler Web处理程序
type WebHandler struct {
	LoggerName string
	Path       string
	Script     string
	Method     string
	Encoding   string
	Handler    func(*Context)
}

//recover 异常处理函数
func (h WebHandler) recover(log logger.ILogger) {
	if r := recover(); r != nil {
		log.Fatal(r, string(debug.Stack()))
	}
}

//call 请求处理回调函数
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
