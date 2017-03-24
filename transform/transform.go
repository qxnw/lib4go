package transform

import (
	"fmt"
	"regexp"
	"sync"
)

//Transform 翻译组件
type Transform struct {
	data  map[string]string
	mutex sync.Mutex
}

//New 创建翻译组件
func New() *Transform {
	return &Transform{data: make(map[string]string)}
}

//NewMaps 根据map创建组件
func NewMaps(d map[string]interface{}) *Transform {
	current := make(map[string]string)
	for k, v := range d {
		current[fmt.Sprintf("@%s", k)] = fmt.Sprint(v)
	}
	return &Transform{data: current}
}

//Set 设置变量的值
func (d *Transform) Set(k string, v string) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.data[fmt.Sprintf("@%s", k)] = v
}

//Get 获取变量的值
func (d *Transform) Get(k string) string {
	return d.data[fmt.Sprintf("@%s", k)]
}

//Translate 翻译带有@变量的字符串
func (d *Transform) Translate(format string) string {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	brackets, _ := regexp.Compile(`\{@\w+\}`)
	result := brackets.ReplaceAllStringFunc(format, func(s string) string {
		return d.data[s[1:len(s)-1]]
	})
	word, _ := regexp.Compile(`@\w+`)
	result = word.ReplaceAllStringFunc(result, func(s string) string {
		return d.data[s]
	})
	return result
}
