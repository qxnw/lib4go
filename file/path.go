package file

import (
	"os"
	"path/filepath"
)

//Exists 检查文件或路径是否存在
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

//GetAbs 获取文件绝对路径
func GetAbs(path string) string {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return ""
	}
	return absPath
}

//CreateFile 根据文件路径(相对或绝对路径)创建文件，如果文件所在的文件夹不存在则自动创建
func CreateFile(path string) (f *os.File, err error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return
	}

	dir := filepath.Dir(absPath)
	err = os.MkdirAll(dir, 0777)
	if err != nil {
		return
	}
	return os.OpenFile(absPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
}
