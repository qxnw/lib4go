package encoding

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"
)

func TestGetReader(t *testing.T) {
	input := "你好"
	charset := "utf-8"
	data, err := ioutil.ReadAll(GetReader(input, charset))
	if err != nil {
		t.Errorf("测试失败：%v", err)
	}
	if !strings.EqualFold(string(data), input) {
		t.Errorf("GetReader fail %s to %s", input, string(data))
	}

	charset = ""
	data, err = ioutil.ReadAll(GetReader(input, charset))
	if bytes.EqualFold(data, []byte(input)) || data == nil {
		t.Error("GetReader fail")
	}
}

func TestCovert(t *testing.T) {
	input := "你好"
	charset := "utf-8"
	except := "你好"
	actual := Convert([]byte(input), charset)
	if !strings.EqualFold(actual, except) {
		t.Errorf("GetReader fail %s to %s", input, actual)
	}

	charset = ""
	data, err := ioutil.ReadAll(GetReader(input, charset))
	if err != nil {
		t.Errorf("测试失败:%v", err)
	}
	actual = Convert([]byte(input), charset)
	if !strings.EqualFold(actual, string(data)) {
		t.Errorf("GetReader fail %s to %s", string(data), actual)
	}
}

func TestUnicodeEncode(t *testing.T) {
	input := "你好"
	except := "\\u4f60\\u597d"
	actual := UnicodeEncode(input)
	if !strings.EqualFold(except, actual) {
		t.Errorf("UnicodeEncode fail %s to %s", except, actual)
	}

	input = "hello world"
	except = "hello world"
	actual = UnicodeEncode(input)
	if !strings.EqualFold(except, actual) {
		t.Errorf("UnicodeEncode fail %s to %s", except, actual)
	}

	input = "!@#!"
	except = "!@#!"
	actual = UnicodeEncode(input)
	if !strings.EqualFold(except, actual) {
		t.Errorf("UnicodeEncode fail %s to %s", except, actual)
	}
}

func TestUnicodeDecode(t *testing.T) {
	input := "\\u4f60\\u597d"
	except := "你好"
	actual := UnicodeDecode(input)
	if !strings.EqualFold(except, actual) {
		t.Errorf("UnicodeDecode fail %s to %s", except, actual)
	}

	input = "\\u0068\\u0065\\u006c\\u006c\\u006f\\u0020\u0077\\u006f\\u0072\\u006c\\u0064"
	except = "hello world"
	actual = UnicodeDecode(input)
	if !strings.EqualFold(except, actual) {
		t.Errorf("UnicodeDecode fail %s to %s", except, actual)
	}

	input = "!@#!"
	except = "!@#!"
	actual = UnicodeDecode(input)
	if !strings.EqualFold(except, actual) {
		t.Errorf("UnicodeDecode fail %s to %s", except, actual)
	}
}
