package net

import (
	"fmt"
	"net/url"
	"strings"
)

//QueryStringToMap 将URL查询字符串中的参数转换成map
func QueryStringToMap(urlQuery string) (result map[string]interface{}, err error) {
	index := strings.IndexAny(urlQuery, "?")
	if index == -1 || index >= len(urlQuery)-1 {
		return nil, nil
	}

	values, err := url.ParseQuery(urlQuery[index+1:])
	if err != nil {
		return nil, fmt.Errorf("url ParseQuery fail: %v", err)
	}
	result = make(map[string]interface{})
	for k, v := range values {
		if len(v) == 1 {
			result[k] = v[0]
		} else {
			result[k] = v
		}
	}

	return
}
