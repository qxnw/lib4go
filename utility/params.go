package utility

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

// GetParamsMap 获取url链接上的数据放到一个map中，通过&分割
func GetParamsMap(urlQuery string) (result map[string]interface{}, err error) {
	/*add by champly 2016年11月18日15:59:36*/
	index := strings.IndexAny(urlQuery, "?")
	if index == -1 || index >= len(urlQuery) {
		return nil, nil
	}
	/*end*/

	values, err := url.ParseQuery(urlQuery[index+1:])
	if err != nil {
		return
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

// GetParams 获取url链接上的数据，并且json格式化，通过&分割
func GetParams(urlQuery string) (res string, err error) {
	result, err := GetParamsMap(urlQuery)
	buffer, err := json.Marshal(&result)
	if err != nil {
		return "", fmt.Errorf("json format fail:%v", err)
	}
	return string(buffer), nil
}
