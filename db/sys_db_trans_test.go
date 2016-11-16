package db

import "testing"

func TestDBTRansQuery(t *testing.T) {
	obj, err := NewSysDB("oracle", "oc_common/123456@orcl136", 2, 2)
	if obj == nil || err != nil {
		t.Error("创建数据库连接失败:", err)
	}

	dbTrans, err := obj.Begin()
	if dbTrans == nil || err != nil {
		t.Errorf("创建数据库事务失败:%v", err)
	}

	// 正常测试
	sql := "select * from oc_user_info where user_id = :1"
	args := []interface{}{"2223"}
	dataRows, colus, err := dbTrans.Query(sql, args...)
	if err != nil {
		t.Errorf("执行%s失败：%v", sql, err)
	}
	if dataRows == nil {
		t.Errorf("执行%s失败", sql)
	}
	if dataRows[0][colus[0]] != "2223" {
		t.Errorf("执行%s失败", sql)
	}

	// 数据库连接串错误测试
	obj, err = NewSysDB("oracle", "", 2, 2)
	if obj != nil || err == nil {
		t.Error("创建数据库连接失败:", err)
	}
	if obj != nil {
		sql = "select * from oc_user_info where user_id = :1"
		args = []interface{}{"2223"}
		dataRows, colus, err = dbTrans.Query(sql, args...)
		if err != nil {
			t.Errorf("执行%s失败：%v", sql, err)
		}
		if dataRows == nil {
			t.Errorf("执行%s失败", sql)
		}
		if dataRows[0][colus[0]] != "2223" {
			t.Errorf("执行%s失败", sql)
		}
	}

	// 数据库连接串错误测试
	obj, err = NewSysDB("", "oc_common/123456@orcl136", 2, 2)
	if obj != nil || err == nil {
		t.Error("创建数据库连接失败:", err)
	}
	if obj != nil {
		sql = "select * from oc_user_info where user_id = :1"
		args = []interface{}{"2223"}
		dataRows, colus, err = dbTrans.Query(sql, args...)
		if err != nil {
			t.Errorf("执行%s失败：%v", sql, err)
		}
		if dataRows == nil {
			t.Errorf("执行%s失败", sql)
		}
		if dataRows[0][colus[0]] != "2223" {
			t.Errorf("执行%s失败", sql)
		}
	}

	// sql错误
	obj, err = NewSysDB("oracle", "oc_common/123456@orcl136", 2, 2)
	if obj == nil || err != nil {
		t.Error("创建数据库连接失败:", err)
	}

	sql = "selects * from dual where 1 = :1"
	args = []interface{}{"1"}
	dataRows, colus, err = dbTrans.Query(sql, args...)
	if err == nil {
		t.Errorf("执行%s失败", sql)
	}

	// sql错误
	obj, err = NewSysDB("oracle", "oc_common/123456@orcl136", 2, 2)
	if obj == nil || err != nil {
		t.Error("创建数据库连接失败:", err)
	}

	sql = "select * from user_id where 1 = :1"
	args = []interface{}{"1"}
	dataRows, colus, err = dbTrans.Query(sql, args...)
	if err == nil {
		t.Errorf("执行%s失败", sql)
	}
}

func TestDBTRansExecute(t *testing.T) {
	obj, err := NewSysDB("oracle", "oc_common/123456@orcl136", 2, 2)
	if obj == nil || err != nil {
		t.Error("创建数据库连接失败:", err)
	}

	dbTrans, err := obj.Begin()
	if dbTrans == nil || err != nil {
		t.Errorf("创建数据库事务失败:%v", err)
	}

	// 正常测试
	sql := "update oc_user_info t set t.traffic_wallet = t.traffic_wallet + 0 where t.user_id = :1"
	args := []interface{}{"2223"}
	row, err := dbTrans.Execute(sql, args...)
	if err != nil {
		t.Errorf("执行%s失败：%v", sql, err)
	}
	if int(row) != 1 {
		t.Errorf("执行%s失败", sql)
	}

	// 数据库连接串错误测试
	obj, err = NewSysDB("oracle", "", 2, 2)
	if obj != nil || err == nil {
		t.Error("创建数据库连接失败:", err)
	}
	if obj != nil {
		sql = "update oc_user_info t set t.traffic_wallet = t.traffic_wallet + 0 where t.user_id = :1"
		args = []interface{}{"2223"}
		row, err = dbTrans.Execute(sql, args...)
		if err != nil {
			t.Errorf("执行%s失败：%v", sql, err)
		}
		if int(row) == 1 {
			t.Errorf("执行%s失败", sql)
		}
	}

	// 数据库连接串错误测试
	obj, err = NewSysDB("", "oc_common/123456@orcl136", 2, 2)
	if obj != nil || err == nil {
		t.Error("创建数据库连接失败:", err)
	}
	if obj != nil {
		sql = "update oc_user_info t set t.traffic_wallet = t.traffic_wallet + 0 where t.user_id = :1"
		args = []interface{}{"2223"}
		row, err = dbTrans.Execute(sql, args...)
		if err != nil {
			t.Errorf("执行%s失败：%v", sql, err)
		}
		if int(row) == 1 {
			t.Errorf("执行%s失败", sql)
		}
	}

	// sql错误
	obj, err = NewSysDB("oracle", "oc_common/123456@orcl136", 2, 2)
	if err != nil {
		t.Error("创建数据库连接失败:", err)
	}
	if obj != nil {
		sql = "updates oc_user_info t set t.traffic_wallet = t.traffic_wallet + 0 where t.user_id = :1"
		args = []interface{}{"2223"}
		row, err = dbTrans.Execute(sql, args...)
		if err == nil {
			t.Errorf("测试失败")
		}
	}

	// sql错误
	obj, err = NewSysDB("oracle", "oc_common/123456@orcl136", 2, 2)
	if err != nil {
		t.Error("创建数据库连接失败:", err)
	}
	if obj != nil {
		sql = "update oc_user_infos t set t.traffic_wallet = t.traffic_wallet + 0 where t.user_id = :1"
		args = []interface{}{"2223"}
		row, err = dbTrans.Execute(sql, args...)
		if err == nil {
			t.Errorf("测试失败")
		}
	}
}

func TestDBTransRollback(t *testing.T) {
	// 正常测试
	obj, err := NewSysDB("oracle", "oc_sales/123456@orcl136", 2, 2)
	if obj == nil || err != nil {
		t.Error("创建数据库连接失败:", err)
	}

	dbTrans, err := obj.Begin()
	if dbTrans == nil || err != nil {
		t.Errorf("创建数据库事务失败:%v", err)
	}
	err = dbTrans.Rollback()
	if err != nil {
		t.Error("回滚数据库事务失败")
	}

	// // 数据库连接串错误
	// obj, err = NewSysDB("oracle", "", 2, 2)
	// if obj != nil || err == nil {
	// 	t.Error("创建数据库连接失败:", err)
	// }

	// err = dbTrans.Rollback()
	// if err != nil {
	// 	t.Error("回滚数据库事务失败")
	// }
}

func TestDBTransCommit(t *testing.T) {
	obj, err := NewSysDB("oracle", "oc_sales/123456@orcl136", 2, 2)
	if obj == nil || err != nil {
		t.Error("创建数据库连接失败:", err)
	}

	dbTrans, err := obj.Begin()
	if dbTrans == nil || err != nil {
		t.Errorf("创建数据库事务失败:%v", err)
	}

	err = dbTrans.Commit()
	if err != nil {
		t.Error("提交数据库事务失败")
	}

	// // 数据库连接串错误
	// obj, err = NewSysDB("oracle", "", 2, 2)
	// if obj != nil || err == nil {
	// 	t.Error("创建数据库连接失败:", err)
	// }

	// err = dbTrans.Commit()
	// if err != nil {
	// 	t.Error("回滚数据库事务失败")
	// }
}