package utility

import (
	"strings"
	"testing"
)

func TestGetSessionID(t *testing.T) {
	sessionID := GetSessionID()
	if sessionID == "" || len(sessionID) != 8 {
		t.Error("获取SessionID失败")
	}
}

func TestGetGUID(t *testing.T) {
	guid := GetGUID()
	if guid == "" {
		t.Error("生成GUID失败")
	}
}

func TestGetLocalIPAddress(t *testing.T) {
	ip := GetLocalIPAddress(nil...)
	if ip == "" {
		t.Error("获取IP地址失败")
	}

	ip = GetLocalIPAddress("192")
	if !strings.HasPrefix(ip, "192") {
		t.Error("获取IP地址失败")
	}

	ip = GetLocalIPAddress("127")
	if !strings.HasPrefix(ip, "127") {
		t.Error("获取IP地址失败")
	}

	ip = GetLocalIPAddress("$%")
	if !strings.HasPrefix(ip, "127") {
		t.Error("获取IP地址失败")
	}
}

func TestEscape(t *testing.T) {
	input := "\\u0026sdfaqr\\\\u0026"
	r := Escape(input)
	if r == input || strings.Contains(r, "\\u0026") {
		t.Error("替换特殊字符失败")
	}

	input = "\\u003csdfaqr\\u003e"
	r = Escape(input)
	if r == input || strings.Contains(r, "\\u003e") || strings.Contains(r, "\\u003c") {
		t.Error("替换特殊字符失败")
	}
}

func TestGetExcPath(t *testing.T) {
	path := GetExcPath(nil...)
	if path != "" {
		t.Error("测试失败")
	}

	path = GetExcPath("./home")
	if path == "" {
		t.Error("获取路径失败")
	}

	path = GetExcPath("/home")
	if path != "/home" {
		t.Error("获取路径失败")
	}
}

// func TestClone(t *testing.T) {
// 	bob := "123"
// 	bob2, err := Clone(bob)
// 	if err != nil {
// 		t.Error("复制失败")
// 	}
// 	if !strings.EqualFold(bob2.(string), bob) {
// 		t.Errorf("复制失败 %s to %s", bob, bob2.(string))
// 	}
// }

func TestGetMin(t *testing.T) {
	a, b := 1, 2
	if a != GetMin(a, b) {
		t.Error("GetMin测试失败")
	}
}

func TestGetMax(t *testing.T) {
	a, b := 1, 2
	if b != GetMax(a, b) {
		t.Error("GetMax测试失败")
	}
}

func TestGetMax2(t *testing.T) {
	a, b, c := 0, 1, 2
	if b != GetMax2(a, b, c) {
		t.Error("GetMax2测试失败")
	}
	if c != GetMax2(b, a, c) {
		t.Error("GetMax2测试失败")
	}
}

func TestCloneMap(t *testing.T) {
	input := make(map[string]interface{})
	input["Name"] = "Bob"
	input["Age"] = 18

	output := CloneMap(input)
	if input["Name"] != output["Name"].(string) || input["Age"] != output["Age"].(int) {
		t.Error("CloneMap测试失败")
	}
}

func TestMerge(t *testing.T) {
	current := make(map[string]interface{})
	current["Name"] = "Bob"
	current["Age"] = 18

	input := make(map[string]interface{})
	input["Name"] = "Tim"
	input["Addr"] = "China"

	Merge(current, input)
	if current["Name"].(string) != input["Name"] || current["Age"].(int) != 18 || current["Addr"].(string) != input["Addr"] {
		t.Error("Merge测试失败")
	}
}

func TestDecodeData(t *testing.T) {
	array := []byte("你好")
	actual, err := DecodeData("UTF-8", array)
	if err != nil {
		t.Error("DecodeData测试失败")
	}
	if !strings.EqualFold(actual, "你好") {
		t.Errorf("DecodeData测试失败:%s", actual)
	}

	actual, err = DecodeData("GBK", array)
	if err != nil {
		t.Error("DecodeData测试失败")
	}
	if actual == "" {
		t.Errorf("DecodeData测试失败:%s", actual)
	}

	actual, err = DecodeData("GB2312", array)
	if err != nil {
		t.Error("DecodeData测试失败")
	}
	if actual == "" {
		t.Errorf("DecodeData测试失败:%s", actual)
	}
}
