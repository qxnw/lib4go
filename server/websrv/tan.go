// Copyright 2015 The WebServer Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package websrv

import (
	ctx "context"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

func Version() string {
	return "0.5.2.1214"
}

type WebServer struct {
	server *http.Server
	Router
	handlers   []Handler
	logger     Logger
	ErrHandler Handler
	ctxPool    sync.Pool
	respPool   sync.Pool
}

var (
	ClassicHandlers = []Handler{
		Logging(),
		Recovery(false),
		Compresses([]string{}),
		Static(StaticOptions{Prefix: "public"}),
		Return(),
		Param(),
		Contexts(),
	}
)

func (t *WebServer) Logger() Logger {
	return t.logger
}

func (t *WebServer) Get(url string, c interface{}, middlewares ...Handler) {
	t.Route([]string{"GET", "HEAD:Get"}, url, c, middlewares...)
}

func (t *WebServer) Post(url string, c interface{}, middlewares ...Handler) {
	t.Route([]string{"POST"}, url, c, middlewares...)
}

func (t *WebServer) Head(url string, c interface{}, middlewares ...Handler) {
	t.Route([]string{"HEAD"}, url, c, middlewares...)
}

func (t *WebServer) Options(url string, c interface{}, middlewares ...Handler) {
	t.Route([]string{"OPTIONS"}, url, c, middlewares...)
}

func (t *WebServer) Trace(url string, c interface{}, middlewares ...Handler) {
	t.Route([]string{"TRACE"}, url, c, middlewares...)
}

func (t *WebServer) Patch(url string, c interface{}, middlewares ...Handler) {
	t.Route([]string{"PATCH"}, url, c, middlewares...)
}

func (t *WebServer) Delete(url string, c interface{}, middlewares ...Handler) {
	t.Route([]string{"DELETE"}, url, c, middlewares...)
}

func (t *WebServer) Put(url string, c interface{}, middlewares ...Handler) {
	t.Route([]string{"PUT"}, url, c, middlewares...)
}

func (t *WebServer) Any(url string, c interface{}, middlewares ...Handler) {
	t.Route(SupportMethods, url, c, middlewares...)
	t.Route([]string{"HEAD:Get"}, url, c, middlewares...)
}

func (t *WebServer) Use(handlers ...Handler) {
	t.handlers = append(t.handlers, handlers...)
}

func GetAddress(args ...interface{}) string {
	var host string
	var port int

	if len(args) == 1 {
		switch arg := args[0].(type) {
		case string:
			addrs := strings.Split(args[0].(string), ":")
			if len(addrs) == 1 {
				host = addrs[0]
			} else if len(addrs) >= 2 {
				host = addrs[0]
				_port, _ := strconv.ParseInt(addrs[1], 10, 0)
				port = int(_port)
			}
		case int:
			port = arg
		}
	} else if len(args) >= 2 {
		if arg, ok := args[0].(string); ok {
			host = arg
		}
		if arg, ok := args[1].(int); ok {
			port = arg
		}
	}

	if host_ := os.Getenv("HOST"); len(host_) != 0 {
		host = host_
	} else if len(host) == 0 {
		host = "0.0.0.0"
	}

	if port_, _ := strconv.ParseInt(os.Getenv("PORT"), 10, 32); port_ != 0 {
		port = int(port_)
	} else if port == 0 {
		port = 8000
	}

	addr := host + ":" + strconv.FormatInt(int64(port), 10)

	return addr
}

// Run the http server. Listening on os.GetEnv("PORT") or 8000 by default.
func (t *WebServer) Run(args ...interface{}) {
	addr := GetAddress(args...)
	t.logger.Info("Listening on http://" + addr)
	t.server = &http.Server{Addr: addr, Handler: t}
	err := t.server.ListenAndServe()
	if err != nil {
		t.logger.Error(err)
	}
}

// Close close the http server
func (t *WebServer) Close(timeout time.Duration) {
	xt, _ := ctx.WithTimeout(ctx.Background(), timeout)
	t.server.Shutdown(xt)
}

func (t *WebServer) RunTLS(certFile, keyFile string, args ...interface{}) {
	addr := GetAddress(args...)
	t.logger.Info("Listening on https://" + addr)

	err := http.ListenAndServeTLS(addr, certFile, keyFile, t)
	if err != nil {
		t.logger.Error(err)
	}
}

type HandlerFunc func(ctx *Context)

func (h HandlerFunc) Handle(ctx *Context) {
	h(ctx)
}

func WrapBefore(handler http.Handler) HandlerFunc {
	return func(ctx *Context) {
		handler.ServeHTTP(ctx.ResponseWriter, ctx.Req())

		ctx.Next()
	}
}

func WrapAfter(handler http.Handler) HandlerFunc {
	return func(ctx *Context) {
		ctx.Next()

		handler.ServeHTTP(ctx.ResponseWriter, ctx.Req())
	}
}

func (t *WebServer) UseHandler(handler http.Handler) {

	t.Use(WrapBefore(handler))
}
func (t *WebServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	resp := t.respPool.Get().(*responseWriter)
	resp.reset(w)

	ctx := t.ctxPool.Get().(*Context)
	ctx.tan = t
	ctx.reset(req, resp)

	ctx.invoke()
	// if there is no logging or error handle, so the last written check.
	if !ctx.Written() {
		p := req.URL.Path
		if len(req.URL.RawQuery) > 0 {
			p = p + "?" + req.URL.RawQuery
		}

		if ctx.Route() != nil {
			if ctx.Result == nil {
				ctx.WriteString("")
				t.logger.Info(req.Method, ctx.Status(), p)
				t.ctxPool.Put(ctx)
				t.respPool.Put(resp)
				return
			}
			panic("result should be handler before")
		}

		if ctx.Result == nil {
			ctx.Result = NotFound()
		}

		ctx.HandleError()

		t.logger.Error(req.Method, ctx.Status(), p)
	}

	t.ctxPool.Put(ctx)
	t.respPool.Put(resp)
}

func NewWithLog(logger Logger, handlers ...Handler) *WebServer {
	tan := &WebServer{
		Router:     NewRouter(),
		logger:     logger,
		handlers:   make([]Handler, 0),
		ErrHandler: Errors(),
	}

	tan.ctxPool.New = func() interface{} {
		return &Context{
			tan:    tan,
			Logger: tan.logger,
		}
	}

	tan.respPool.New = func() interface{} {
		return &responseWriter{}
	}

	tan.Use(handlers...)

	return tan
}

func New(handlers ...Handler) *WebServer {
	return NewWithLog(NewLogger(os.Stdout), handlers...)
}

func Classic(l ...Logger) *WebServer {
	var logger Logger
	if len(l) == 0 {
		logger = NewLogger(os.Stdout)
	} else {
		logger = l[0]
	}

	return NewWithLog(
		logger,
		ClassicHandlers...,
	)
}
