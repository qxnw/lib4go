package db

import (
	"fmt"
	"sync"
	"testing"
)

// BenchmarkSysDbTest 测试数据库连接的并发性
func BenchmarkSysDbTest(b *testing.B) {
	sql := "select * from test where id = :1"
	args := []interface{}{"1"}
	wg := &sync.WaitGroup{}
	n := b.N
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			obj, err := NewSysDB("oracle", dbConnectStr, 1, 2)
			if obj == nil || err != nil {
				b.Error("创建数据库连接失败:", err)
			}
			dataRows, colus, err := obj.Query(sql, args...)
			if err != nil {
				b.Errorf("test fail %v", err)
			}
			obj.Close()
			fmt.Printf("%+v\t%+v\n", dataRows, colus)
			wg.Done()
		}()
	}

	wg.Wait()

	// obj, err := NewSysDB("oracle", "oc_test/123456@orcl136", 2, 5)
	// if obj == nil || err != nil {
	// 	b.Error("创建数据库连接失败:", err)
	// }

	// sql := "select * from test where id = :1"
	// // sql := "update test t set t.money = :1 where t.id = 1"

	// // num := 0

	// wg := &sync.WaitGroup{}

	// n := 20
	// fmt.Println(n)
	// wg.Add(n)

	// for i := 0; i < n; i++ {
	// 	go func(i int) {
	// 		args := []interface{}{i}
	// 		dataRows, colus, err := obj.Query(sql, args...)
	// 		// data, err := obj.Execute(sql, args)
	// 		if err != nil {
	// 			b.Errorf("执行%s失败：%v", sql, err)
	// 		}

	// 		// fmt.Println(err, result)
	// 		// if dataRows == nil {
	// 		// 	b.Errorf("执行%s失败", sql)
	// 		// }
	// 		if len(dataRows) > 0 {
	// 			fmt.Printf("%+v\t%+v\n", dataRows, colus)
	// 		}
	// 		// if dataRows[0][colus[0]] != "1" {
	// 		// 	b.Errorf("执行%s失败", sql)
	// 		// }
	// 		wg.Done()
	// 	}(i)
	// }
	// wg.Wait()
}
