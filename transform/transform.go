package transform

import (
	"fmt"
	"regexp"
	"sync"
)

type ITransformGetter interface {
	Set(string, string)
	Get(string) string
}
type transformData map[string]string

func (t transformData) Get(key string) string {
	if v, ok := t[key]; ok {
		return fmt.Sprintf("%v", v)
	}
	return ""
}
func (t transformData) Set(key string, value string) {
	t[key] = value
}

//Transform 翻译组件
type Transform struct {
	Data  ITransformGetter
	mutex sync.Mutex
}

//New 创建翻译组件
func New() *Transform {
	var data transformData = make(map[string]string)
	return &Transform{Data: data}
}

//NewGetter getter
func NewGetter(t ITransformGetter) *Transform {
	return &Transform{Data: t}
}

//NewMaps 根据map创建组件
func NewMaps(d map[string]interface{}) *Transform {
	var data transformData = make(map[string]string)
	for k, v := range d {
		data[k] = fmt.Sprint(v)
	}
	return &Transform{Data: data}
}

//NewMap create by map
func NewMap(d map[string]string) *Transform {
	var data transformData = make(map[string]string)
	for k, v := range d {
		data[k] = fmt.Sprint(v)
	}
	return &Transform{Data: data}
}

//Set 设置变量的值
func (d *Transform) Set(k string, v string) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.Data.Set(k, v)
}

//Get 获取变量的值
func (d *Transform) Get(k string) string {
	return d.Data.Get(k)
}

//Translate 翻译带有@变量的字符串
func (d *Transform) Translate(format string) string {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	brackets, _ := regexp.Compile(`\{@\w+\}`)
	result := brackets.ReplaceAllStringFunc(format, func(s string) string {
		return d.Data.Get(s[2 : len(s)-1])
	})
	word, _ := regexp.Compile(`@\w+`)
	result = word.ReplaceAllStringFunc(result, func(s string) string {
		return d.Data.Get(s[1:])
	})
	return result
}
