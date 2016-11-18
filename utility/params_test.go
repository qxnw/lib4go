package utility

import "testing"

func TestGetParamsMap(t *testing.T) {
	url := "http://geek.csdn.net/news/detail/124352?locationNum=2&fps=1"
	data, err := GetParamsMap(url)
	if err != nil {
		t.Errorf("GetParamsMap fail: %v", err)
	}
	if len(data) != 2 || data["locationNum"] != "2" || data["fps"] != "1" {
		t.Errorf("GetParamsMap fail:%s", data)
	}

	url = "http://geek.csdn.net/news/detail/124352??locationNum=2&fps=1"
	data, err = GetParamsMap(url)
	if err != nil {
		t.Errorf("GetParamsMap fail: %v", err)
	}
	if len(data) != 2 || data["?locationNum"] != "2" || data["fps"] != "1" {
		t.Errorf("GetParamsMap fail:%s", data)
	}

	url = "http://geek.csdn.net/news/detail/124352?"
	data, err = GetParamsMap(url)
	if err != nil {
		t.Errorf("GetParamsMap fail: %v", err)
	}
	if len(data) != 0 {
		t.Errorf("GetParamsMap fail:%s", data)
	}

	url = "http://geek.csdn.net/news/detail/124352?locationNum"
	data, err = GetParamsMap(url)
	if err != nil {
		t.Errorf("GetParamsMap fail: %v", err)
	}
	if len(data) != 1 || data["locationNum"] != "" {
		t.Errorf("GetParamsMap fail:%s", data)
	}

	url = "http://geek.csdn.net/news/detail/12435"
	data, err = GetParamsMap(url)
	if err != nil {
		t.Errorf("GetParamsMap fail: %v", err)
	}
	if len(data) != 0 {
		t.Errorf("GetParamsMap fail:%s", data)
	}

	url = "asdfaqrew"
	data, err = GetParamsMap(url)
	if err != nil {
		t.Errorf("GetParamsMap fail: %v", err)
	}
	if len(data) != 0 {
		t.Errorf("GetParamsMap fail:%s", data)
	}
}

func TestGetParams(t *testing.T) {
	url := "http://geek.csdn.net/news/detail/124352?locationNum=2&fps=1"
	data, err := GetParams(url)
	if err != nil {
		t.Errorf("GetParams fail: %v", err)
	}
	if data == "{}" {
		t.Errorf("GetParams fail:%s", data)
	}

	url = "http://geek.csdn.net/news/detail/124352??locationNum=2&fps=1"
	data, err = GetParams(url)
	if err != nil {
		t.Errorf("GetParams fail: %v", err)
	}
	if data == "{}" {
		t.Errorf("GetParams fail:%s", data)
	}

	url = "http://geek.csdn.net/news/detail/124352?"
	data, err = GetParams(url)
	if err != nil {
		t.Errorf("GetParams fail: %v", err)
	}
	if data != "{}" {
		t.Errorf("GetParams fail:%s", data)
	}

	url = "http://geek.csdn.net/news/detail/124352"
	data, err = GetParams(url)
	if err != nil {
		t.Errorf("GetParams fail: %v", err)
	}
	if data != "{}" {
		t.Errorf("GetParams fail:%s", data)
	}

	url = "1safsdf"
	data, err = GetParams(url)
	if err != nil {
		t.Errorf("GetParams fail: %v", err)
	}
	if data != "{}" {
		t.Errorf("GetParams fail:%s", data)
	}
}
