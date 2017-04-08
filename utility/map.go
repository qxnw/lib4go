package utility

import (
	"encoding/json"
	"net/url"
)

//GetMapWithQuery 将URL参数转换为map
func GetMapWithQuery(query string) (r map[string]string, err error) {
	values, err := url.ParseQuery(query)
	if err != nil {
		return
	}
	r = make(map[string]string)
	for k, v := range values {
		if len(v) >= 0 {
			r[k] = v[0]
		}
	}
	return
}

//GetJSONWithQuery 将URL参数转换为JSON
func GetJSONWithQuery(query string) (res string, err error) {
	result, err := GetMapWithQuery(query)
	buffer, err := json.Marshal(&result)
	if err != nil {
		return
	}
	return string(buffer), nil
}
