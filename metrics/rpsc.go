package metrics

import (
	"sync/atomic"
	"time"
)

//RPS 基于HashedWheelTimer算法的计数器，过期自动淘汰
type RPSC struct {
	total      int64
	length     int64
	slots      []int64
	lastTicker int64
	counter    int64
}

//NewRPSC 构建计数器
func NewRPSC(length int64, total int64) (w *RPSC) {
	w = &RPSC{length: length, total: total}
	w.slots = make([]int64, w.total, w.total)
	for i := int64(0); i < w.total; i++ {
		w.slots[i] = 0
	}
	return w
}

//Mark 添加新值
func (r *RPSC) Mark(new int64) {
	r.mark(new, time.Now().Unix())
}

//mark 记录上次执行时间，超过时间间隔则清除counter
//每一跳需清除中间秒数
func (r *RPSC) mark(new int64, currentStep int64) {
	lastStep := r.lastTicker
	current := currentStep % r.total
	atomic.AddInt64(&r.counter, -r.clear(lastStep, currentStep)) //6, 8(clear,1,2,7,8)
	atomic.AddInt64(&r.counter, new)
	atomic.AddInt64(&r.slots[current], new)
	r.lastTicker = currentStep
}

func (r *RPSC) clear(l, n int64) (clearCounter int64) { //1-5:1,10:1,10 //2:1,3:1
	clearCounter = int64(0)
	if l == n {
		return
	}
	//清空时间中间差
	if n-l >= r.length {
		for i := int64(0); i < r.total; i++ {
			clearCounter += atomic.SwapInt64(&r.slots[i], 0)
		}
		return clearCounter
	}

	right := n % r.total                   //0,3
	l1 := (right - r.length + 1) % r.total //5,4
	left := l1 % r.total
	if l1 < 0 {
		left = (l1 + r.total) % r.total
	}
	if right > left {
		for i := int64(0); i < left; i++ { //0,1,2,3,4,5
			clearCounter += atomic.SwapInt64(&r.slots[i], 0)
		}
		for i := right; i < r.total; i++ { //1,
			clearCounter += atomic.SwapInt64(&r.slots[i], 0)
		}
		return clearCounter
	}
	for i := right; i < left; i++ { //0,1,2,3,4,5 //3,4
		clearCounter += atomic.SwapInt64(&r.slots[i], 0)
	}
	return clearCounter

}
