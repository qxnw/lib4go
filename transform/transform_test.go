package transform

import (
	"fmt"
	"strings"
	"testing"
)

type TestType struct {
	name string
	age  int
}

// TestNew 测试创建一个空的翻译组件
func TestNew(t *testing.T) {
	// 创建一个翻译组件，判断结果
	transform := New()
	if len(transform.data) != 0 {
		t.Errorf("test fail actual : %d", len(transform.data))
	}
}

// TestNewMaps 测试根据map创建一个翻译组件
func TestNewMaps(t *testing.T) {
	input := map[string]map[string]interface{}{
		"value1":   map[string]interface{}{"key1": "value1"},
		"!@#!%^%$": map[string]interface{}{"key2": "!@#!%^%$"},
		"123":      map[string]interface{}{"key3": 123},
		"[228 189 160 229 165 189]": map[string]interface{}{"key4": []byte("你好")},
		"{Tom 12}":                  map[string]interface{}{"key5": TestType{name: "Tom", age: 12}},
		"<nil>":                     map[string]interface{}{"key6": nil},
	}

	for except, actual := range input {
		transform := NewMaps(actual)
		for k := range transform.data {
			if !strings.EqualFold(transform.data[k], except) {
				t.Errorf("test fail actual:%s, except:%s", transform.data[k], except)
			}
		}

	}
}

// TestSet 测试设置变量的值
func TestSet(t *testing.T) {
	// 构建一个空的翻译组件
	transform := New()

	input := map[string]map[string]string{
		"value1":   map[string]string{"key1": "value1"},
		"!@#!%^%$": map[string]string{"key2": "!@#!%^%$"},
		"123":      map[string]string{"key3": fmt.Sprintf("%d", 123)},
		"你好":       map[string]string{"key4": string([]byte("你好"))},
	}

	for except, actual := range input {
		for k, v := range actual {
			transform.Set(k, v)
			if !strings.EqualFold(transform.data["@"+k], except) {
				t.Errorf("test fail actual:%s, except:%s", transform.data["@"+k], except)
			}
		}

	}
}

// TestGet 测试获取组件中的值
func TestGet(t *testing.T) {
	// 构建一个空的翻译组件
	transform := New()

	input := map[string]map[string]string{
		"value1":   map[string]string{"key1": "value1"},
		"!@#!%^%$": map[string]string{"key2": "!@#!%^%$"},
		"123":      map[string]string{"key3": fmt.Sprintf("%d", 123)},
		"你好":       map[string]string{"key4": string([]byte("你好"))},
	}

	for except, actual := range input {
		for k, v := range actual {
			transform.Set(k, v)
			if !strings.EqualFold(transform.Get(k), except) {
				t.Errorf("test fail actual:%s, except:%s", transform.Get(k), except)
			}
		}
	}

	// 获取一个不存在的值
	except := ""
	actual := transform.Get("s2vs1!$")
	if !strings.EqualFold(actual, except) {
		t.Errorf("test fail actual:%s, except:%s", actual, except)
	}
}

// TestTranslate 测试翻译带有@变量的字符串
func TestTranslate(t *testing.T) {
	// 构建一个空的翻译组件
	transform := New()

	datas := map[string]map[string]string{
		"value1":   map[string]string{"key1": "value1"},
		"!@#!%^%$": map[string]string{"key2": "!@#!%^%$"},
		"123":      map[string]string{"key3": fmt.Sprintf("%d", 123)},
		"你好":       map[string]string{"key4": string([]byte("你好"))},
	}

	for _, data := range datas {
		for k, v := range data {
			transform.Set(k, v)
		}
	}

	input := "asdfa@key1, {@key2 }, {@key3}, @@key4@ @@ key4 @key11 {@key4"
	except := "asdfavalue1, {!@#!%^%$ }, 123, @你好@ @@ key4  {你好"
	actual := transform.Translate(input)
	if !strings.EqualFold(actual, except) {
		t.Errorf("test fail, actual : %s, except : %s", actual, except)
	}
}
