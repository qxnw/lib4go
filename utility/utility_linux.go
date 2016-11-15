package utility

import (
	"os"
	"strings"
)

func getExecRoot() (path string, err error) {
	f, err := os.Readlink("proc/self/exe")
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
