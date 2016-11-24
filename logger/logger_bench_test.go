package logger

import (
	"testing"
	"time"
)

func BenchLogger(b *testing.B) {
	for i := 0; i < b.N; i++ {
		go func(i int) {
			l := New("test", "test1")
			for j := 0; j < 5; j++ {
				l.Debug("当前数量:", i)
				time.Sleep(time.Second)
			}
		}(i)
	}

	time.Sleep(time.Second * 15)

	Close()
}
