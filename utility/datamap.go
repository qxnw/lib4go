package utility

import (
	"fmt"
	"regexp"
	"sync"
)

type DataMap struct {
	data  map[string]string
	mutex sync.Mutex
}

func NewDataMap() DataMap {
	return DataMap{data: make(map[string]string)}
}
func NewDataMaps(d map[string]interface{}) DataMap {
	current := make(map[string]string)
	for k, v := range d {
		current[fmt.Sprintf("@%s", k)] = fmt.Sprint(v)
	}
	return DataMap{data: current}
}

//Add 添加变量
func (d DataMap) Set(k string, v string) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.data[fmt.Sprintf("@%s", k)] = v
}
func (d DataMap) Get(k string) string {
	return d.data[fmt.Sprintf("@%s", k)]
}

//Merge merge new map from current
func (d DataMap) Merge(n DataMap) DataMap {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	nmap := NewDataMap()
	MergeStringMap(d.data, nmap.data)
	MergeStringMap(n.data, nmap.data)
	return nmap
}

//Copy Copy the current map to another
func (d DataMap) Copy() DataMap {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	nmap := NewDataMap()
	MergeStringMap(d.data, nmap.data)
	return nmap
}

//Translate 翻译带有@变量的字符串
func (d DataMap) Translate(format string) string {
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

func MergeMaps(source map[string]interface{}, targets []map[string]interface{}) []map[string]interface{} {
	for k, v := range source {
		for _, target := range targets {
			target[k] = v
		}
	}
	return targets
}
func MergeMap(source map[string]interface{}, target map[string]interface{}) map[string]interface{} {
	for k, v := range source {
		target[k] = v
	}
	return target
}
func MergeStringMap(source map[string]string, target map[string]string) map[string]string {
	for k, v := range source {
		target[k] = v
	}
	return target
}
