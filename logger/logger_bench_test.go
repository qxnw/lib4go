package logger

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/qxnw/lib4go/file"

	"io/ioutil"

	"strings"
)

// 测试输出到文件
func BenchmarkLoggerToFile(b *testing.B) {
	// for i := 0; i < b.N; i++ {
	// 	go func(i int) {
	// 		l := New("test", "test1")
	// 		for j := 0; j < 5; j++ {
	// 			l.Debug("当前数量:", i)
	// 			time.Sleep(time.Second)
	// 		}
	// 	}(i)
	// }

	// time.Sleep(time.Second * 15)

	// Close()

	// TimeWriteToFile = 1
	// manager = newLoggerManager()
	// log := New("name1", "name2", "name3")
	// for i := 0; i < b.N; i++ {
	// 	// go func(i int) {
	// 	// 	for j := 0; j < 10; j++ {
	// 	// 		log.Debugf("携程编号：%d", i)
	// 	// 	}
	// 	// }(i)
	// 	for j := 0; j < 10; j++ {
	// 		log.Debugf("携程编号：%d", i)
	// 	}
	// }

	// // time.Sleep(time.Second * 1)
	// Close()
	// TimeWriteToFile = time.Second

	// 把数据写入文件
	totalAccount := 10000 * 5
	lk := sync.WaitGroup{}
	ch := make(chan int, totalAccount)
	name := "ABC"

	log := New(name)

	doWriteToFile := func(ch chan int, lk *sync.WaitGroup) {
	START:
		for {
			select {
			case l, ok := <-ch:
				if ok {
					log.Debug(l)
					log.Info(l)
					log.Error(l)
				} else {
					break START
				}
				lk.Done()
			}
		}
	}

	for i := 0; i < 100; i++ {
		go doWriteToFile(ch, &lk)
	}

	for i := 0; i < totalAccount; i++ {
		lk.Add(1)
		ch <- i
	}
	close(ch)
	lk.Wait()

	// time.Sleep(time.Second * 1)

	Close()

	// 开始读取文件
	path := fmt.Sprintf("../logs/%s/%d%d%d.log", name, time.Now().Year(), time.Now().Month(), time.Now().Day())
	filePath := file.GetAbs(path)
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		b.Errorf("test fail : %v", err)
	}
	count := len(strings.Split(string(data), "\n"))
	if count != totalAccount*3+6 {
		b.Errorf("test fail, actual:%d, except:%d", count, totalAccount*3+6)
	}

	// // 删除文件
	// os.Remove(filePath)

}
