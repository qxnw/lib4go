package utility

import "os"

func GetExecRoot() (path string, err error) {
	return os.Getwd()
}
