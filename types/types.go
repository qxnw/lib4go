package types

import (
	"fmt"
)

//DecodeString 判断变量的值与指定相等时设置为另一个值，否则使用原值
func DecodeString(def interface{}, a interface{}, b interface{}, e ...interface{}) string {
	values := make([]interface{}, 0, len(e)+2)
	values = append(values, a)
	values = append(values, b)
	values = append(values, e...)

	for i := 0; i < len(values); i = i + 2 {
		if def == values[i] {
			return fmt.Sprint(values[i+1])
		}
	}
	if len(values)%2 == 1 {
		return fmt.Sprint(values[len(values)-1])
	}
	return ""
}
