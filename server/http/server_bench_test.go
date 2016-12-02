package http

import (
	"testing"
	"time"

	"github.com/qxnw/lib4go/net/http"
)

func BenchmarkServer(b *testing.B) {
	totalAccount := 0
	// var totalTime time.Duration

	address := "localhost:8080"
	loggerName := "test"
	path := "/api/test"
	script := "test.lua"
	method := "get"
	encoding := "utf-8"
	handler := func(context *Context) {
		totalAccount++
		// totalTime += context.PassTime()
	}
	handlers := NewHandler(loggerName, path, script, method, encoding, handler)
	server := NewWebServer(address, loggerName, handlers)

	// 开启server服务
	go server.Serve()
	time.Sleep(time.Second)

	client := http.NewHTTPClient()

	total := b.N
	for i := 0; i < total; i++ {
		client.Get("http://localhost:8080/api/test", method)
	}

	// 执行完校验结果
	if total != totalAccount {
		b.Errorf("test fail actual:%d, except:%d", totalAccount, total)
	}

	// b.Logf("总共处理次数：%d\t平均处理时间：%dns", totalAccount, int32(totalTime)/int32(totalAccount))

	// 关闭服务
	server.Stop()
}
