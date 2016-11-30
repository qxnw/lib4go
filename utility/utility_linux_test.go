package utility

import (
	"strings"
	"testing"
)

func TestLinuxGetExecRoot(t *testing.T) {
	path, err := getExecRoot()
	if err != nil {
		t.Errorf("test fail : %v", err)
	}
	if !strings.EqualFold(path, "") {
		t.Error("test fail")
	}
}
