package utility

import (
	"testing"
)

// TestGetGUID 测试生成的Guid是否重复
func TestGetGUID(t *testing.T) {
	totalAccount := 10000 * 1000
	data := map[string]int{}

	for i := 0; i < totalAccount; i++ {
		key := GetGUID()
		data[key] = i
	}

	if len(data) != totalAccount {
		t.Errorf("test fail, totalAccount:%d, actual:%d", totalAccount, len(data))
	}
}
