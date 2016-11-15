package utility

import (
	"encoding/json"
	"net/url"
)

func GetParamsMap(urlQuery string) (result map[string]interface{}, err error) {
	values, err := url.ParseQuery(urlQuery)
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

func GetParams(urlQuery string) (res string, err error) {
	result, err := GetParamsMap(urlQuery)
	buffer, err := json.Marshal(&result)
	if err != nil {
		return
	}
	return string(buffer), nil
}
