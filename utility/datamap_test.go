package utility

import (
	"testing"
)

func TestNewDataMap(t *testing.T) {
	dataMap := NewDataMap()
	if len(dataMap.data) != 0 {
		t.Errorf("test fail : %d", len(dataMap.data))
	}
}

func TestNewDataMaps(t *testing.T) {
	inputs := map[string]interface{}{
		`string`: map[string]interface{}{"key": "string"},
	}
}
