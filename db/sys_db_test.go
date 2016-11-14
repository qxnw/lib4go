package db

import "testing"

func TestNewSysDB(t *testing.T) {
	input := []SysDB{
		{"sqlite", "oc_common/123456@orcl136", nil, 2, 2},
		{"oracle", "oc_common/123456@orcl136", nil, 2, 2},
		{"oracle", "oc_common/123456@orcl136", nil, 4, 2},
		{"oracle", "oc_common/123456@orcl136", nil, -2, 2},
		{"oracle", "oc_common/123456@orcl136", nil, 2, -2},
		// {"1231", "oc_common/123456@orcl136", nil, 1, 2},
		// {"oracle", "", nil, 2, 2},
		// {"oracle", "^&&^@#@", nil, 2, 2},
		// {"oracle", "oc_common/123456@orcl136", nil, 2, 2},
	}
	for _, data := range input {
		obj, err := NewSysDB(data.provider, data.connString, data.maxIdle, data.maxOpen)
		if obj == nil || err != nil {
			t.Errorf("测试失败:%v", err)
		}
	}

	// 测试不支持的类型
	_, err := NewSysDB("1231", "oc_common/123456@orcl136", 2, 2)
	if err == nil {
		t.Error("测试失败")
	}

	// // 连接串错误
	// _, err = NewSysDB("oracle", "", 2, 2)
	// if err == nil {
	// 	t.Error("测试失败")
	// }

	// // 连接串错误
	// _, err = NewSysDB("oracle", "^&&^@#@", 2, 2)
	// if err == nil {
	// 	t.Error("测试失败")
	// }

	// 数据库配置错误
	_, err = NewSysDB("oracle", "oc_common/123456@orcl136", -2, 2)
	if err != nil {
		t.Error("测试失败")
	}

	// 数据库配置错误
	_, err = NewSysDB("oracle", "oc_common/123456@orcl136", 2, -2)
	if err != nil {
		t.Error("测试失败")
	}
}

type testQueryResult struct {
	args []map[string]string
	data string
	row  int
}

func TestQuery(t *testing.T) {
	// 正常流程
	obj, err := NewSysDB("oracle", "oc_common/123456@orcl136", 2, 2)
	if obj == nil || err != nil {
		t.Error("创建数据库连接失败:", err)
	}

	sql := "select * from oc_user_info where user_id = :1"
	args := []interface{}{"2223"}
	dataRows, colus, err := obj.Query(sql, args...)
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
	if obj != nil || err != nil {
		t.Error("创建数据库连接失败:", err)
	}
	if obj != nil {
		sql = "select * from oc_user_info where user_id = :1"
		args = []interface{}{"2223"}
		dataRows, colus, err = obj.Query(sql, args...)
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
	dataRows, colus, err = obj.Query(sql, args...)
	if err == nil {
		t.Errorf("执行%s失败", sql)
	}
}

func TestExecute(t *testing.T) {
	obj, err := NewSysDB("oracle", "oc_common/123456@orcl136", 2, 2)
	if obj == nil || err != nil {
		t.Error("创建数据库连接失败:", err)
	}

	sql := "update oc_user_info t set t.traffic_wallet = t.traffic_wallet + 0 where t.user_id = :1"
	args := []interface{}{"2223"}
	row, err := obj.Execute(sql, args...)
	if err != nil {
		t.Errorf("执行%s失败：%v", sql, err)
	}
	if int(row) != 1 {
		t.Errorf("执行%s失败", sql)
	}

	// 数据库连接串错误测试
	obj, err = NewSysDB("oracle", "", 2, 2)
	if obj != nil || err != nil {
		t.Error("创建数据库连接失败:", err)
	}
	if obj != nil {
		sql = "update oc_user_info t set t.traffic_wallet = t.traffic_wallet + 0 where t.user_id = :1"
		args = []interface{}{"2223"}
		row, err = obj.Execute(sql, args...)
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
		row, err = obj.Execute(sql, args...)
		if err == nil {
			t.Errorf("测试失败")
		}
	}
}

func TestBegin(t *testing.T) {
	obj, err := NewSysDB("oracle", "oc_common/123456@orcl136", 2, 2)
	if obj == nil || err != nil {
		t.Error("创建数据库连接失败:", err)
	}

	dbTrans, err := obj.Begin()
	if dbTrans == nil || err != nil {
		t.Errorf("创建数据库事务失败:%v", err)
	}
}

func TestClose(t *testing.T) {
	obj, err := NewSysDB("oracle", "oc_common/123456@orcl136", 2, 2)
	if obj == nil || err != nil {
		t.Error("创建数据库连接失败:", err)
	}

	obj.Close()
}
