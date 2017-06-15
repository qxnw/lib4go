package metrics

import (
	"sync/atomic"
	"time"
)

//RPS 基于HashedWheelTimer算法的计数器，过期自动淘汰
type RPSC struct {
	total      int32
	length     int32
	slots      []int32
	lastTicker int64
	counter    int32
}

//NewRPS 构建计数器
func NewRPSC(length int32, total int32) (w *RPSC) {
	w = &RPSC{length: length, total: total}
	w.slots = make([]int32, w.total, w.total)
	for i := int32(0); i < w.total; i++ {
		w.slots[i] = 0
	}
	return w
}

//Mark 添加新值
func (r *RPSC) Mark(new int32) {
	r.mark(new, time.Now().Unix())
}

//mark 记录上次执行时间，超过时间间隔则清除counter
//每一跳需清除中间秒数
func (r *RPSC) mark(new int32, currentStep int64) {
	lastStep := r.lastTicker
	current := int32(currentStep) % r.total
	atomic.AddInt32(&r.counter, -r.clear(lastStep, currentStep)) //6, 8(clear,1,2,7,8)
	atomic.AddInt32(&r.counter, new)
	atomic.AddInt32(&r.slots[current], new)
	r.lastTicker = currentStep
}

func (r *RPSC) clear(l, n int64) (clearCounter int32) { //1-5:1,10:1,10 //2:1,3:1
	clearCounter = int32(0)
	//清空时间中间差
	if int32(n-l) >= r.length {
		for i := int32(0); i < r.total; i++ {
			clearCounter += atomic.SwapInt32(&r.slots[i], 0)
		}
		return clearCounter
	}

	right := int32(n % int64(r.total))             //0,3
	l1 := int32(int32(right)-r.length+1) % r.total //5,4
	left := l1 % r.total
	if l1 < 0 {
		left = (l1 + r.total) % r.total
	}
	if right > left {
		for i := int32(0); i < left; i++ { //0,1,2,3,4,5
			clearCounter += atomic.SwapInt32(&r.slots[i], 0)
		}
		for i := right; i < r.total; i++ { //1,
			clearCounter += atomic.SwapInt32(&r.slots[i], 0)
		}
		return clearCounter
	}
	for i := right; i < left; i++ { //0,1,2,3,4,5 //3,4
		clearCounter += atomic.SwapInt32(&r.slots[i], 0)
	}
	return clearCounter

}
