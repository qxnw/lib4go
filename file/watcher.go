package file

import (
	"os"
	"path/filepath"
	"time"

	"sync"

	"github.com/qxnw/lib4go/concurrent/cmap"
)

//DirWatcher 文件夹监控
type DirWatcher struct {
	callback func()
	files    cmap.ConcurrentMap
	lastTime time.Time
	timeSpan time.Duration
	done     bool
	mu       sync.Mutex
}

//NewDirWatcher 构建脚本监控文件
func NewDirWatcher(callback func(), timeSpan time.Duration) *DirWatcher {
	w := &DirWatcher{callback: callback, lastTime: time.Now(), timeSpan: timeSpan}
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
	for {
		select {
		case <-time.After(time.Second):
			if w.done {
				return
			}
		case <-time.After(w.timeSpan):
			if w.done {
				return
			}
			if w.checkChange() {
				w.callback()
			}
		}
	}
}

//checkChange 检查文件夹最后修改时间
func (w *DirWatcher) checkChange() bool {
	w.mu.Lock()
	defer w.mu.Unlock()
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

//Close 关闭服务
func (w *DirWatcher) Close() {
	w.done = true
}
