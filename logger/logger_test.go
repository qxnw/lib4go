package logger

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	logger := New(nil...)
	fmt.Println(logger)
}
