package utility

import "testing"

func TestGetExecRoot(t *testing.T) {
	path := getExecRoot(nil...)
	if path != "" {
		t.Error("测试失败")
	}

	path = getExecRoot("./home")
	if path == "" {
		t.Error("获取路径失败")
	}

	path = getExecRoot("/home")
	if path != "/home" {
		t.Error("获取路径失败")
	}
}
