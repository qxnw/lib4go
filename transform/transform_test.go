package transform

import (
	"testing"
)

func TestNew(t *testing.T) {
	// 创建一个翻译组件，判断结果
	transform := New()
	if len(transform.data) != 0 {
		t.Errorf("test fail actual : %d", len(transform.data))
	}
}

func TestNewMaps(t *testing.T) {
	input := map[string]interface{}{
	// {"test1": "test1"},
	// {"test2": "test2"},
	}
}
