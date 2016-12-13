package db

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// BenchmarkSysDbTest 测试数据库连接的并发性
func TestSysDbTest(b *testing.T) {
	sql := "select * from test where id = :1"
	args := []interface{}{"1"}
	wg := &sync.WaitGroup{}
	n := 50
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			obj, err := NewSysDB("oracle", dbConnectStr, 10, 20)
			if obj == nil || err != nil {
				b.Error("创建数据库连接失败:", err)
			}
			_, _, err = obj.Query(sql, args...)
			if err != nil {
				b.Errorf("test fail %v", err)
			}
			obj.Close()
			// fmt.Printf("%+v\t%+v\n", dataRows, colus)
			wg.Done()
		}()
	}
	wg.Wait()
}

func TestSysDbTest2(b *testing.T) {
	sql := "select * from test where id = :1"

	args := []interface{}{"1"}
	wg := &sync.WaitGroup{}
	times := time.Second * 0

	for j := 0; j < 100; j++ {
		obj, err := NewSysDB("oracle", dbConnectStr, 1, 1)
		if obj == nil || err != nil {
			b.Error("创建数据库连接失败:", err)
		}
		n := 300
		// fmt.Println(n)
		wg.Add(n)
		start := time.Now()
		for i := 0; i < n; i++ {
			go func() {
				// obj.Print()
				_, _, err := obj.Query(sql, args...)
				if err != nil {
					b.Errorf("test fail %v", err)
				}
				// obj.Print()
				// obj.Close()
				// fmt.Printf("%+v\t%+v\n", dataRows, colus)
				wg.Done()
			}()
		}
		wg.Wait()
		// fmt.Println("总共耗时：", time.Now().Sub(start))
		times += time.Now().Sub(start)
	}

	fmt.Println(times / 100)
}
