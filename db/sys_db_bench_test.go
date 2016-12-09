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
	n := 1000
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			obj, err := NewSysDB("oracle", dbConnectStr, 1, 2)
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

func BenchmarkSysDbTest2(b *testing.B) {
	// sql := "select * from test where id = :1"
	sql := `select t.order_no,
       t.business_type,
       t.recharge_account,
       t.face,
       t.create_time,
       t.standard,
       t.delivery_status,
       t.order_status,
       p.product_name PRODUCTNAME,
       c.channel_name CHANNELNAME
  from (select rid
          from (select rid, rownum linenum
                  from (select rowid rid
                          from zb_order_info t
                         where t.create_time >=
                               to_date('2016/12/7 0:00:00',
                                       'yyyy-mm-dd hh24:mi:ss')
                           and t.create_time <=
                               to_date('2016/12/9 0:00:00',
                                       'yyyy-mm-dd hh24:mi:ss')
                        and 1 = :1
                        )
 where rownum <= 1 * 10)
 where linenum > 10 * (1 - 1)) TAB1
 inner join zb_order_info t on t.rowid = TAB1.rid
 inner join zb_up_channel_product p on t.up_product_no = p.product_no
 inner join zb_up_channel_info c on t.up_channel_no = c.channel_no
`
	args := []interface{}{"1"}
	wg := &sync.WaitGroup{}
	obj, err := NewSysDB("oracle", "zb_sales/123456@ORCL136", 1, 2)
	if obj == nil || err != nil {
		b.Error("创建数据库连接失败:", err)
	}
	n := 500
	fmt.Println(n)
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			// obj.Print()
			_, _, err := obj.Query(sql, args...)
			if err != nil {
				b.Errorf("test fail %v", err)
			}
			obj.Print()
			// obj.Close()
			// fmt.Printf("%+v\t%+v\n", dataRows, colus)
			wg.Done()
		}()
	}
	wg.Wait()
}
