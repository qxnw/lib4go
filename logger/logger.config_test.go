package logger

import "testing"

func TestRead(t *testing.T) {
	// 配置文件不存在
	loggerPath = "../conf/no_ars.logger.json"
	_, err := read()
	t.Log(err)
	if err == nil {
		t.Error("test fail")
		return
	}

	// 无法读取配置文件,文件权限为000
	loggerPath = "../conf/without_x_ars.logger.json"
	_, err = read()
	t.Log(err)
	if err == nil {
		t.Error("test fail")
		return
	}

	// 配置文件有误
	loggerPath = "../conf/err_ars.logger.json"
	_, err = read()
	t.Log(err)
	if err == nil {
		t.Error("test fail")
		return
	}

	// 正常读取配置文件
	loggerPath = "../conf/ars.logger.json"
	_, err = read()
	t.Log(err)
	if err != nil {
		t.Errorf("读取配置文件失败：%v", err)
		return
	}
}

func TestWriteToFile(t *testing.T) {
	// 正常读取配置文件
	loggerPath = "../conf/ars.logger.json"
	appenders, err := read()
	if err != nil {
		t.Errorf("读取配置文件失败：%v", err)
		return
	}

	// 读取配置文件失败【没有权限】
	err = writeToFile("/root/fail.log", appenders)
	if err == nil {
		t.Error("test fail")
	}

	// 配置文件不存在
	err = writeToFile("../conf/no_ars.logger.json", appenders)
	if err != nil {
		t.Errorf("test fail：%v", err)
	}

	// 正常写配置文件
	err = writeToFile("../logs/test.log", appenders)
	if err != nil {
		t.Errorf("test fail：%v", err)
	}
}
