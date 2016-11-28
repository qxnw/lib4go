/*
	关于Http Server的文档：
		http://www.cnblogs.com/yjf512/archive/2012/08/22/2650873.html
*/

package http

import (
	"net"
	"net/http"

	"github.com/qxnw/lib4go/logger"
)

//WebServer WEB服务
type WebServer struct {
	routes     []*WebHandler
	address    string
	loggerName string
	Log        logger.ILogger
	l          net.Listener
}

//NewWebServer 创建WebServer服务
func NewWebServer(address string, loggerName string, handlers ...*WebHandler) (server *WebServer) {
	server = &WebServer{routes: handlers, address: address, loggerName: loggerName}
	server.Log = logger.Get(loggerName)
	return
}

//Serve 启动WEB服务器
func (s *WebServer) Serve() (err error) {
	mux := http.NewServeMux()
	for _, handler := range s.routes {
		mux.HandleFunc(handler.Path, handler.call)

	}
	l, err := net.Listen("tcp", s.address)
	if err != nil {
		return
	}
	s.l = l
	err = http.Serve(s.l, mux)
	return
}

//Stop 停止服务器
func (s *WebServer) Stop() {
	if s.l != nil {
		s.l.Close()
	}
}
