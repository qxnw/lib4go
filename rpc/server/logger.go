package server

import (
	"io"

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
	l := log.New(out, "[RpcServer] ", log.Ldefault())
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

	}
}
