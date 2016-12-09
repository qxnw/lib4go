package file

import (
	"os"
	"path/filepath"
	"time"

	"github.com/qxnw/lib4go/concurrent/cmap"
)

//DirWatcher 文件夹监控
type DirWatcher struct {
	callback func()
	files    cmap.ConcurrentMap
	lastTime time.Time
}

//NewDirWatcher 构建脚本监控文件
func NewDirWatcher(callback func()) *DirWatcher {
	w := &DirWatcher{callback: callback, lastTime: time.Now()}
	w.files = cmap.New()
	go w.watch()
	return w
}

//Append 添加监控文件
func (w *DirWatcher) Append(path string) (err error) {
	dir := filepath.Dir(path)
	w.files.SetIfAbsent(dir, dir)
	return nil
}

func (w *DirWatcher) watch() {
	tk := time.NewTicker(time.Second * 10)
	for {
		select {
		case <-tk.C:
			if w.checkChange() {
				w.callback()
			}
		}
	}
}

//checkChange 检查文件夹最后修改时间
func (w *DirWatcher) checkChange() bool {
	change := false
	w.files.IterCb(func(path string, v interface{}) bool {
		fileinfo, err := os.Stat(path)
		if err != nil {
			return !change
		}
		if fileinfo.ModTime().Sub(w.lastTime) > 0 {
			w.lastTime = time.Now()
			change = true
			return !change
		}
		return !change
	})
	return change
}
