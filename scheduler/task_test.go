package scheduler

import (
	"testing"
	"time"
)

func BenchmarkTask(b *testing.B) {
	count := b.N
	real := 0
	notifyChan := make(chan int, count)
	for i := 0; i < count; i++ {
		AddTask("0/1 * * * * ?", NewTask(notifyChan, func(v interface{}) {
			nchan := v.(chan int)
			nchan <- 1
		}))
	}
	Start()
	ticker := time.NewTicker(time.Second * 5)
Loop:
	for {
		select {
		case <-notifyChan:
			real++
			if real == count {
				break Loop
			}
		case <-ticker.C:
			break Loop
		}
	}
	if real != count {
		b.Error("请求超时", real)
	}
}
