/*
{
    "address":"http://192.168.101.183:8086",
     "db":"influxdb_test",
     "user":"admin",
     "password":"123456",
     "row":"test_table,name=@name,address=@address age=@age"
}
*/
package influxdb

import "testing"

func TestNewJSON(t *testing.T) {
	// 参数正确
	config := `{
    "address":"http://192.168.0.185:8086",
    "db":"influxdb_test",
    "user":"admin",
    "password":"123456",
    "row":"test_table,name=@name,address=@address age=@age"
	}`
	_, err := NewJSON(config)
	if err != nil {
		t.Errorf("create influxdb context fail:%v", err)
	}

	// 缺少参数
	config = `{
	"address":"http://192.168.0.185:8086",
	"user":"admin",
	"password":"123456",
	"row":"test_table,name=@name,address=@address age=@age"
	}`
	_, err = NewJSON(config)
	if err == nil {
		t.Error("test fail")
	}

	// 缺少参数
	config = `{
	"address":"http://192.168.0.185:8086",
	"db":"influxdb_test",
	"row":"test_table,name=@name,address=@address age=@age"
	}`
	_, err = NewJSON(config)
	if err != nil {
		t.Errorf("create influxdb context fail:%v", err)
	}

	// 缺少参数
	config = `{
	"address":"http://192.168.0.185:8086",
	"db":"influxdb_test",
	"user":"admin",
	"password":"123456"
	}`
	_, err = NewJSON(config)
	if err == nil {
		t.Error("test fail")
	}

	// 配置有误
	config = `
	"address":"http://192.168.0.185:8086",
	"db":"influxdb_test",
	"user":"admin",
	"password":"123456"
	}`
	_, err = NewJSON(config)
	if err == nil {
		t.Error("test fail")
	}
}

func TestSaveString(t *testing.T) {
	// 正常流程
	config := `{
    "address":"http://192.168.0.185:8086",
    "db":"influxdb_test",
    "user":"admin",
    "password":"123456",
    "row":"test_table,name=@name,address=@address age=@age"
	}`
	i, err := NewJSON(config)
	if err != nil {
		t.Errorf("create influxdb context fail:%v", err)
	}

	rows := `[{"name":"champly","address":"china","age":"18"}]`
	err = i.SaveString(rows)
	if err != nil {
		t.Errorf("SaveString fail:%v", err)
	}

	// 参数错误
	config = `{
    "address":"http://192.168.0.185:8086",
    "db":"influxdb_test",
    "user":"admin",
    "password":"123456",
    "row":"test_table,name=@name,address=@address age=@age"
	}`
	i, err = NewJSON(config)
	if err != nil {
		t.Errorf("create influxdb context fail:%v", err)
	}

	rows = `[{"name":"champly","address":"china"}]`
	err = i.SaveString(rows)
	if err == nil {
		t.Error("test SaveString fail")
	}

	// influxdb地址错误
	config = `{
    "address":"http://192.168.0.186:8086",
    "db":"influxdb_test",
    "user":"admin",
    "password":"123456",
    "row":"test_table,name=@name,address=@address age=@age"
	}`
	i, err = NewJSON(config)
	if err != nil {
		t.Errorf("create influxdb context fail:%v", err)
	}

	rows = `[{"name":"champly","address":"china","age":"18"}]`
	err = i.SaveString(rows)
	if err == nil {
		t.Error("test SaveString fail")
	}

	// 数据库错误
	config = `{
    "address":"http://192.168.0.185:8086",
    "db":"influxdb_test_err",
    "user":"admin",
    "password":"123456",
    "row":"test_table,name=@name,address=@address age=@age"
	}`
	i, err = NewJSON(config)
	if err != nil {
		t.Errorf("create influxdb context fail:%v", err)
	}

	rows = `[{"name":"champly","address":"china","age":"18"}]`
	err = i.SaveString(rows)
	if err == nil {
		t.Error("test SaveString fail")
	}

	// 参数错误
	config = `{
    "address":"http://192.168.0.185:8086",
    "db":"influxdb_test",
    "user":"admin",
    "password":"123456",
    "row":"test_table,name=@name,address=@address age=@age"
	}`
	i, err = NewJSON(config)
	if err != nil {
		t.Errorf("create influxdb context fail:%v", err)
	}

	rows = `["name":"champly","address":"china","age":"18"}]`
	err = i.SaveString(rows)
	if err == nil {
		t.Error("test SaveString fail")
	}
}
