package utility

import (
	"os"
	"path/filepath"
)

// getExecRoot 获取当前可执行文件的绝对路径的目录
func getExecRoot() (path string, err error) {
	f, err := filepath.Abs(os.Args[0])
	if err != nil {
		return
	}
	/*change by champly 2016年11月23日09:32:06*/
	// index := strings.LastIndex(f, "/")
	// if index > -1 {
	// 	path = f[:index]
	// 	return
	// }
	/*end*/
	path = f
	return
}
