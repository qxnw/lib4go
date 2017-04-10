package transform

import "testing"
import "github.com/qxnw/lib4go/ut"

type TestType struct {
	name string
	age  int
}
type getter struct {
	data map[string]string
}

func (g getter) Get(key string) string {
	return g.data[key]
}
func (g getter) Set(key string, value string) {
	g.data[key] = value
}

// TestNewMaps 测试根据map创建一个翻译组件
func TestNewMaps(t *testing.T) {
	m := NewMaps(map[string]interface{}{
		"id":   "123",
		"name": "colin",
	})
	ut.Expect(t, m.Translate("@id"), "123")
	ut.Expect(t, m.Translate("@name"), "colin")
	ut.Expect(t, m.Translate("@name/@id"), "colin/123")
	ut.Expect(t, m.Translate("{@name}/{@id}"), "colin/123")
	m.Set("age", "100")
	ut.Expect(t, m.Translate("{@name}/{@id}/@age"), "colin/123/100")
	ut.Expect(t, m.Translate("{@name}/{@id}/@age/@age2"), "colin/123/100/")
}

// TestNewMaps 测试根据map创建一个翻译组件
func TestNewGetter(t *testing.T) {
	m := NewGetter(getter{data: map[string]string{"id": "123", "name": "colin"}})
	ut.Expect(t, m.Translate("@id"), "123")
	ut.Expect(t, m.Translate("@name"), "colin")
	ut.Expect(t, m.Translate("@name/@id"), "colin/123")
	ut.Expect(t, m.Translate("{@name}/{@id}"), "colin/123")
	m.Set("age", "100")
	ut.Expect(t, m.Translate("{@name}/{@id}/@age"), "colin/123/100")
	ut.Expect(t, m.Translate("{@name}/{@id}/@age/@age2"), "colin/123/100/")
}
