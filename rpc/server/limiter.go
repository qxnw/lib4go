package server

import "github.com/qxnw/lib4go/metrics"

type limiter struct {
	data map[string]float64
}

func (m limiter) Handle(ctx *Context) {
	service := ctx.Req().Service
	if count, ok := m.data["*"]; ok {
		conterName := metrics.MakeName(ctx.server.serverName+".limit_", "service", "*")
		meter := metrics.GetOrRegisterMeter(conterName, metrics.DefaultRegistry)
		if meter.Rate1() >= count {
			ctx.Forbidden()
			ctx.server.logger.Error("service:%s 超过总流程限制(%d)", service, count)
			return
		}
		meter.Mark(1)
	}
	if count, ok := m.data[service]; ok {
		conterName := metrics.MakeName(ctx.server.serverName+".limit_", "service", service)
		meter := metrics.GetOrRegisterMeter(conterName, metrics.DefaultRegistry)
		if meter.Rate1() >= count {
			ctx.Forbidden()
			ctx.server.logger.Error("service:%s 超过当前服务的流程限制(%d)", service, count)
			return
		}
		meter.Mark(1)
	}

	ctx.Next()
}
