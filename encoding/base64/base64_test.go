package base64

import (
	"bytes"
	"strings"
	"testing"
)

func TestEncodeBytes(t *testing.T) {
	input := "你好"
	except := "5L2g5aW9"
	actual := Encode(input)

	if !strings.EqualFold(except, actual) {
		t.Errorf("Encode %s fail,%s", input, actual)
	}

	buf := []byte(input)
	actual = EncodeBytes(buf)
	if !strings.EqualFold(except, actual) {
		t.Errorf("EncodeBytes %s fail,%s", input, actual)
	}

	input = "5L2g5aW9"
	result := []byte("你好")
	actualBuf, err := DecodeBytes(input)
	if err != nil {
		t.Error("DecodeBytes fail")
		return
	}
	if !bytes.EqualFold(actualBuf, result) {
		t.Error("DecodeBytes fail")
	}

	except = "你好"
	actual, err = Decode(input)
	if err != nil {
		t.Error("Decode fail")
		return
	}
	if !strings.EqualFold(actual, except) {
		t.Error("Decode fail")
	}

	errInput := "!@#!"
	_, err = DecodeBytes(errInput)
	if err == nil {
		t.Error("测试错误")
		return
	}

	_, err = Decode(errInput)
	if err == nil {
		t.Error("测试错误")
		return
	}

}
