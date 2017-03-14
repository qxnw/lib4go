package websrv

import (
	"time"

	"github.com/qxnw/lib4go/metrics"
)

type Metric struct {
	Host     string
	Database string
	username string
	password string
	timeSpan time.Duration
}

func NewMetric(host string, dataBase string, userName string, password string, timeSpan time.Duration) *Metric {
	m := &Metric{Host: host, Database: dataBase, username: userName, password: password, timeSpan: timeSpan}
	go metrics.InfluxDB(metrics.DefaultRegistry, timeSpan, m.Host, m.Database,
		m.username,
		m.password)
	return m
}
func (m *Metric) Execute(context *Context) {
	if action := context.Action(); action != nil {
		if l, ok := action.(LogInterface); ok {
			l.SetLogger(context.Logger)
		}
	}
	context.Next()
}
func (m *Metric) Handle(ctx *Context) {
	start := time.Now()
	p := ctx.Req().URL.Path

	//add meter
	meter := metrics.GetOrRegisterMeter(p, metrics.DefaultRegistry)
	meter.Mark(1)

	//add counter
	counter := metrics.GetOrRegisterCounter(p, metrics.DefaultRegistry)
	counter.Inc(1)

	// add time
	timer := metrics.GetOrRegisterTimer(p, metrics.DefaultRegistry)
	timer.Time(func() { m.Execute(ctx) })

	counter.Dec(1)

	if !ctx.Written() {
		if ctx.Result == nil {
			ctx.Result = NotFound()
		}
		ctx.HandleError()
	}

	statusCode := ctx.Status()

	if statusCode >= 200 && statusCode < 400 {
		ctx.Info(ctx.Req().Method, statusCode, time.Since(start), p)
	} else {
		ctx.Error(ctx.Req().Method, statusCode, time.Since(start), p, ctx.Result)
	}

}
