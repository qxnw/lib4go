package utility

import "os"


// getExecRoot 如果可执行文件是一个链接文件，那么返回的是原文件的目录，如果不是，则返回当前可执行文件的目录
func getExecRoot() (path string, err error) {
	// f, err := os.Readlink("proc/self/exe")
	// if err != nil {
	// 	return
	// }
	// index := strings.LastIndex(f, "/")
	// if index > -1 {
	// 	path = f[:index]
	// 	return
	// }
	// path = f
	// return

	/*change by champly 2016年11月30日09:56:19*/
	f, err := os.Readlink("/proc/self/exe")
	if err != nil {
		return
	}
	path = f
	return
	/*end*/
}
