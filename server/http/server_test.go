package http

import (
	"strings"
	"testing"

	"github.com/qxnw/lib4go/net/http"
)

// TestNewWebServer 测试创建一个webserver服务
func TestNewWebServer(t *testing.T) {
	except := "hello world"

	address := "localhost:8080"
	loggerName := "test"
	path := "/api/test"
	script := "test.lua"
	method := "get"
	encoding := "utf-8"
	handler := func(context *Context) {
		context.Writer.Write([]byte(except))
	}
	handlers := NewHandler(loggerName, path, script, method, encoding, handler)
	server := NewWebServer(address, loggerName, handlers)

	// 开启server服务
	go server.Serve()

	// 通过http请求，访问Server服务，校验数据
	client := http.NewHTTPClient()

	// 正常请求
	actual, status, err := client.Get("http://localhost:8080/api/test", encoding)
	if err != nil {
		t.Errorf("test fail : %v", err)
	}
	if status != 200 {
		t.Errorf("test fail, status is %d", status)
	}
	if !strings.EqualFold(actual, except) {
		t.Errorf("test fail actual : %s, except : %s", actual, except)
	}

	// 通过不同的method
	actual, status, err = client.Post("http://localhost:8080/api/test", encoding)
	if err != nil {
		t.Errorf("test fail : %v", err)
	}
	if status != 404 {
		t.Errorf("test fail, status is %d", status)
	}
	if !strings.EqualFold(actual, "您访问的页面不存在") {
		t.Errorf("test fail actual : %s, except : %s", actual, "您访问的页面不存在")
	}

	// 错误的url
	actual, status, err = client.Post("http://localhost:8080/api/test12", encoding)
	if err != nil {
		t.Errorf("test fail : %v", err)
	}
	if status != 404 {
		t.Errorf("test fail, status is %d", status)
	}
	if !strings.Contains(actual, "404 page not found") {
		t.Errorf("test fail actual : %s, except : %s", actual, "404 page not found")
	}

	// 关闭掉Server服务
	server.Stop()

}

// TestSepcialSituation 测试特殊情况
func TestSepcialSituation(t *testing.T) {
	// 监听一个错误的端口
	except := "hello world"
	address := "localhost:22"
	loggerName := "test"
	path := "/api/test"
	script := "test.lua"
	method := "get"
	encoding := "utf-8"
	handler := func(context *Context) {
		context.Writer.Write([]byte(except))
	}
	handlers := NewHandler(loggerName, path, script, method, encoding, handler)
	server := NewWebServer(address, loggerName, handlers)

	// 开启server服务
	err := server.Serve()
	t.Log(err)
	if err == nil {
		t.Errorf("test fail")
	}

	// 提供所有的method
	address = "localhost:8080"
	method = "*"
	handlers = NewHandler(loggerName, path, script, method, encoding, handler)
	server = NewWebServer(address, loggerName, handlers)

	go server.Serve()

	// 通过http请求，访问Server服务，校验数据
	client := http.NewHTTPClient()

	// 正常请求
	actual, status, err := client.Get("http://localhost:8080/api/test", encoding)
	if err != nil {
		t.Errorf("test fail : %v", err)
	}
	if status != 200 {
		t.Errorf("test fail, status is %d", status)
	}
	if !strings.EqualFold(actual, except) {
		t.Errorf("test fail actual : %s, except : %s", actual, except)
	}

	// 关闭掉Server服务
	server.Stop()

	// 多次关闭服务
	server.Stop()
}
