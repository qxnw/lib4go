package server

import (
	"time"

	"github.com/qxnw/lib4go/concurrent/cmap"
	"github.com/qxnw/lib4go/metrics"
)

type InfluxMetric struct {
	Host     string
	Database string
	username string
	password string
	timeSpan time.Duration
	registry cmap.ConcurrentMap
}

func NewInfluxMetric(host string, dataBase string, userName string, password string, timeSpan time.Duration) *InfluxMetric {
	m := &InfluxMetric{Host: host, Database: dataBase, username: userName, password: password, timeSpan: timeSpan}
	go metrics.InfluxDB(metrics.DefaultRegistry, timeSpan, m.Host, m.Database,
		m.username,
		m.password)
	go metrics.DefaultRegistry.RunHealthchecks()
	return m
}

func (m *InfluxMetric) Execute(context *Context) {
	if action := context.Action(); action != nil {
		if l, ok := action.(LogInterface); ok {
			l.SetLogger(context.Logger)
		}
	}
	context.Next()
}
func (m *InfluxMetric) Handle(ctx *Context) {
	service := ctx.Req().Service
	client := ctx.IP()
	conterName := metrics.MakeName(ctx.server.serverName+".request", metrics.COUNTER, "server", ctx.server.address, "client", client, "service", service)
	timerName := metrics.MakeName(ctx.server.serverName+".request", metrics.TIMER, "server", ctx.server.address, "client", client, "service", service)
	successName := metrics.MakeName(ctx.server.serverName+".success", metrics.METER, "server", ctx.server.address, "client", client, "service", service)
	failedName := metrics.MakeName(ctx.server.serverName+".failed", metrics.METER, "server", ctx.server.address, "client", client, "service", service)

	counter := metrics.GetOrRegisterCounter(conterName, metrics.DefaultRegistry)
	counter.Inc(1)

	metrics.GetOrRegisterTimer(timerName, metrics.DefaultRegistry).Time(func() { m.Execute(ctx) })
	counter.Dec(1)

	if !ctx.Written() {
		if ctx.Result == nil {
			ctx.Result = NotFound()
		}
		ctx.HandleError()
	}

	statusCode := ctx.Writer.Code
	if statusCode >= 200 && statusCode < 400 {
		metrics.GetOrRegisterMeter(successName, metrics.DefaultRegistry).Mark(1)
	} else {
		metrics.GetOrRegisterMeter(failedName, metrics.DefaultRegistry).Mark(1)
	}
}
