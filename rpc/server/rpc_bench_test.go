package server

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/qxnw/grpc4ars/client"
	"github.com/qxnw/grpc4ars/server"
)

func BenchmarkItems(t *testing.B) {
	svr := server.NewServer(func(session string, svs string, data string) (status int, result string, err error) {
		status = 100
		result = svs
		return
	})
	go func() {
		if err := svr.Start(":10165"); err != nil {
			t.Error(err)
		}
	}()

	client := client.NewClient(":10165")
	client.Connect()
	fmt.Println("s:", time.Now())
	s, result, err := client.Request("123455666", "svname", "{}")
	//	fmt.Println(svname)
	if err != nil {
		t.Error(err)
	}
	if s != 100 || result != "svname" {
		t.Error("数据有误")
	}
	fmt.Println("e:", time.Now())

	mu := sync.WaitGroup{}
	for i := 0; i < t.N; i++ {
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
}
