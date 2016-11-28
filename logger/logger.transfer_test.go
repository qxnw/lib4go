package logger

import (
	"strings"
	"testing"
	"time"
)

func TestTransform(t *testing.T) {
	tpls := map[string][]interface{}{
		``:         []interface{}{LogEvent{Level: "Info", Now: time.Now(), Name: "name1", Session: "12345678", Content: "content1", Output: "output1"}, ``},
		`%session`: []interface{}{LogEvent{Level: "Info", Now: time.Now(), Name: "name1", Session: "12345678", Content: "content1", Output: "output1"}, `12345678`},
		`%date`:    []interface{}{LogEvent{Level: "Info", Now: time.Now(), Name: "name1", Session: "12345678", Content: "content1", Output: "output1"}, `20161128`},
		`%year`:    []interface{}{LogEvent{Level: "Info", Now: time.Now(), Name: "name1", Session: "12345678", Content: "content1", Output: "output1"}, `2016`},
		`%mm`:      []interface{}{LogEvent{Level: "Info", Now: time.Now(), Name: "name1", Session: "12345678", Content: "content1", Output: "output1"}, `11`},
		`%dd`:      []interface{}{LogEvent{Level: "Info", Now: time.Now(), Name: "name1", Session: "12345678", Content: "content1", Output: "output1"}, `28`},
		`%level`:   []interface{}{LogEvent{Level: "Info", Now: time.Now(), Name: "name1", Session: "12345678", Content: "content1", Output: "output1"}, `Info`},
		`%l`:       []interface{}{LogEvent{Level: "Info", Now: time.Now(), Name: "name1", Session: "12345678", Content: "content1", Output: "output1"}, `I`},
		`%name`:    []interface{}{LogEvent{Level: "Info", Now: time.Now(), Name: "name1", Session: "12345678", Content: "content1", Output: "output1"}, `name1`},
		`%content`: []interface{}{LogEvent{Level: "Info", Now: time.Now(), Name: "name1", Session: "12345678", Content: "content1", Output: "output1"}, `content1`},
		`%test`:    []interface{}{LogEvent{Level: "Info", Now: time.Now(), Name: "name1", Session: "12345678", Content: "content1", Output: "output1"}, ``},
		`test`:     []interface{}{LogEvent{Level: "Info", Now: time.Now(), Name: "name1", Session: "12345678", Content: "content1", Output: "output1"}, `test`},
	}

	for tpl, except := range tpls {
		actual := transform(tpl, except[0].(LogEvent))
		if !strings.EqualFold(actual, except[1].(string)) {
			t.Errorf("test fail actual：%s\texcept:%s\t", actual, except[1].(string))
		}
	}
}
