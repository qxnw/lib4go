package server

import (
	"sync"
	"testing"

	"google.golang.org/grpc"

	"time"

	"fmt"

	"github.com/qxnw/grpc4ars/client"
	"github.com/qxnw/grpc4ars/server"
	"github.com/qxnw/lib4go/logger"
)

func init() {
	grpc.EnableTracing = false

}

//测试服务正常调用
func TestNormal(t *testing.T) {
	svr := server.NewServer(func(session string, svs string, data string) (status int, result string, err error) {
		status = 100
		result = svs
		return
	})
	go func() {
		if err := svr.Start(":10160"); err != nil {
			t.Error(err)
		}
	}()

	client := client.NewClient(":10160")
	client.Connect()
	s, result, err := client.Request("123455666", "svname", "{}")
	if err != nil {
		t.Error(err)
	}
	if s != 100 || result != "svname" {
		t.Error("数据有误")
	}

	mu := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		mu.Add(1)
		go func(i int) {
			svname := fmt.Sprintf("svs:%d", i)
			s, result, err := client.Request("123455666", svname, "{}")
			//	fmt.Println(svname)
			if err != nil {
				t.Error(err)
			}
			if s != 100 || result != svname {
				t.Error("数据有误")
			}
			mu.Done()

		}(i)
	}
	mu.Wait()
	svr.Close()
	client.Close()

}

//测试服务正常到服务停止，再到服务恢复
func TestReconnect(t *testing.T) {

	//------创建服务器
	svr := server.NewServer(func(session string, svs string, data string) (status int, result string, err error) {
		status = 100
		result = svs
		return
	})
	go func() {
		if err := svr.Start(":10161"); err != nil {
			t.Error(err)
		}
	}()
	time.Sleep(time.Second)
	//--------客户端连接到服务器
	log := logger.Get("test")
	client := client.NewClient(":10161", client.WithHeartbeat(), client.WithCheckTicker(time.Second), client.WithLogger(log))
	client.Connect()
	client.Connect()
	//----------发送请求
	s, result, err := client.Request("123455666", "svname", "{}")
	if err != nil {
		t.Error(err)
	}
	if s != 100 || result != "svname" {
		t.Error("数据有误")
	}
	//---------关闭服务
	svr.Close()
	time.Sleep(time.Second * 3) // 等待端口释放

	//-------发送请求
	s, result, err = client.Request("123455666", "svname12345", "{}")
	if err == nil {
		t.Error("网络连接已关闭应返回失败", s, result)
	}
	//------重新启动服务
	go func() {
		fmt.Println("重新启动服务")
		if err := svr.Start(":10161"); err != nil {
			t.Error(err)
		}
	}()
	time.Sleep(time.Second * 3)
	//-------再次请求服务
	s, result, err = client.Request("123455666", "svname", "{}")
	if err != nil {
		t.Error(err)
	}
	if s != 100 || result != "svname" {
		t.Error("数据有误")
	}
}

//测试客户端由无法连接到服务启动后恢复请求
func TestCantConnect(t *testing.T) {
	log := logger.Get("test")
	//------连接到一个不可达的服务
	client := client.NewClient(":10163", client.WithLogger(log))
	if client.Connect(); client.IsConnect {
		t.Error("应返回无法连接到服务器")
		return
	}

	//------启动服务器
	fmt.Printf("网络连接状态:%v\n", client.IsConnect)
	svr := server.NewServer(func(session string, svs string, data string) (status int, result string, err error) {
		status = 100
		result = svs
		return
	})
	go func() {
		if err := svr.Start(":10163"); err != nil {
			t.Error(err)
		}
	}()
	time.Sleep(time.Second)
	fmt.Printf("网络连接状态:%v\n", client.IsConnect)

	//------客户端自动请求服务
	s, result, err := client.Request("123455666", "svname", "{}")
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("网络连接状态:%v\n", client.IsConnect)
	if !client.IsConnect || s != 100 || result != "svname" {
		t.Error("数据有误")
	}
	svr.Close()
}
