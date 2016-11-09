package tpl

import (
	"fmt"
	"regexp"
)

func isNil(input interface{}) bool {
	return input == nil || fmt.Sprintf("%v", input) == ""
}

//AnalyzeTPL 解析模板内容，并返回解析后的SQL语句，入输入参数
//@表达式，替换为参数化字符如: :1,:2,:3
//#表达式，替换为指定值，值为空时返回NULL
//~表达式，检查值，值为空时返加"",否则返回: , name=value
//&条件表达式，检查值，值为空时返加"",否则返回: and name=value
//|条件表达式，检查值，值为空时返回"", 否则返回: or name=value

func AnalyzeTPL(tpl string, input map[string]interface{}, prefix func() string) (sql string, params []interface{}) {
	params = make([]interface{}, 0)
	word, _ := regexp.Compile(`[@|#|&|~|\||!]\w+`)
	//@变量, 将数据放入params中
	sql = word.ReplaceAllStringFunc(tpl, func(s string) string {
		key := s[1:]
		pre := s[:1]

		value := input[key]
		switch pre {
		case "@":
			if !isNil(value) {
				params = append(params, value)
			} else {
				params = append(params, nil)
			}
			return prefix()
		case "#":
			if !isNil(value) {
				return fmt.Sprintf("%v", value)
			}
			return "NULL"
		case "&":
			if !isNil(value) {
				params = append(params, value)
				return fmt.Sprintf("and %s=%s", key, prefix())
			}
			return ""
		case "|":
			if !isNil(value) {
				params = append(params, value)
				return fmt.Sprintf("or %s=%s", key, prefix())
			}
			return ""
		case "~":
			if !isNil(value) {
				params = append(params, value)
				return fmt.Sprintf(",%s=%s", key, prefix())
			}
			return ""
		default:
			return s
		}
	})
	return
}
