package server

import (
	"io"
	"time"

	"github.com/lunny/log"
)

type Logger interface {
	Debugf(format string, v ...interface{})
	Debug(v ...interface{})
	Infof(format string, v ...interface{})
	Info(v ...interface{})
	Warnf(format string, v ...interface{})
	Warn(v ...interface{})
	Errorf(format string, v ...interface{})
	Error(v ...interface{})
}

func NewLogger(out io.Writer) Logger {
	l := log.New(out, "[rpc server] ", log.Ldefault())
	l.SetOutputLevel(log.Ldebug)
	return l
}

type LogInterface interface {
	SetLogger(Logger)
}

type Log struct {
	Logger
}

func (l *Log) SetLogger(log Logger) {
	l.Logger = log
}

func Logging() HandlerFunc {
	return func(ctx *Context) {
		start := time.Now()
		ctx.server.logger.Info("Started", ctx.Req().Service, "for", ctx.Req().GetArgs()["session"])

		if action := ctx.Action(); action != nil {
			if l, ok := action.(LogInterface); ok {
				l.SetLogger(ctx.Logger)
			}
		}

		ctx.Next()

		if !ctx.Written() {
			if ctx.Result == nil {
				ctx.Result = NotFound()
			}
			ctx.HandleError()
		}

		statusCode := ctx.Writer.Code
		if statusCode >= 200 && statusCode < 400 {
			ctx.server.logger.Info(ctx.Req().Service, statusCode, time.Since(start), ctx.Result)
		} else {
			ctx.server.logger.Error(ctx.Req().Service, statusCode, time.Since(start), ctx.Result)
		}
	}
}
