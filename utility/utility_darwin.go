package utility

import (
	"os"
	"path/filepath"
	"strings"
)

func getExecRoot() (path string, err error) {
	f, err := filepath.Abs(os.Args[0])
	if err != nil {
		return
	}
	index := strings.LastIndex(f, "/")
	if index > -1 {
		path = f[:index]
		return
	}
	path = f
	return
}
