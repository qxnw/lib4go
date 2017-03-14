package websrv

import (
	"time"

	"github.com/qxnw/lib4go/concurrent/cmap"
	"github.com/qxnw/lib4go/metrics"
)

type Metric struct {
	Host     string
	Database string
	username string
	password string
	timeSpan time.Duration
	registry cmap.ConcurrentMap
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
	url := ctx.Req().URL.Path
	client := ctx.IP()
	metricReuqestName := metrics.MakeName(ctx.tan.serverName+".request", "server", ctx.tan.ip, "client", client, "url", url)
	metricSuccessName := metrics.MakeName(ctx.tan.serverName+".request", "server", ctx.tan.ip, "client", client, "url", url)
	metricFailedName := metrics.MakeName(ctx.tan.serverName+".request", "server", ctx.tan.ip, "client", client, "url", url)

	//add meter
	//meter := metrics.GetOrRegisterMeter(metricReuqestName+".meter", metrics.DefaultRegistry)
	//meter.Mark(1)

	//add counter
	counter := metrics.GetOrRegisterCounter("counter."+metricReuqestName, metrics.DefaultRegistry)
	counter.Inc(1)

	// add time
	timer := metrics.GetOrRegisterTimer("timer."+metricReuqestName, metrics.DefaultRegistry)
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
		ctx.Info(ctx.Req().Method, statusCode, time.Since(start), url)
		meter := metrics.GetOrRegisterMeter(metricSuccessName, metrics.DefaultRegistry)
		meter.Mark(1)

	} else {
		ctx.Error(ctx.Req().Method, statusCode, time.Since(start), url, ctx.Result)
		meter := metrics.GetOrRegisterMeter(metricFailedName, metrics.DefaultRegistry)
		meter.Mark(1)
	}

}
