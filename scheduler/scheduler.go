package scheduler

import (
	"sync/atomic"

	"github.com/arsgo/cron"
)

//Scheduler 定时任务调度组件
type Scheduler struct {
	c      *cron.Cron
	status int32
}

//New 创建 定时任务调度组件
func New() *Scheduler {
	return &Scheduler{c: cron.New(), status: 1}
}

//AddTask 添加调度任务
func (s *Scheduler) AddTask(trigger string, task *Task) {
	s.c.AddJob(trigger, task)
}

//Start 启动组件
func (s *Scheduler) Start() {
	if atomic.CompareAndSwapInt32(&s.status, 1, 0) {
		s.c.Start()
	}
}

//Stop 停止组件并清空所有任务
func (s *Scheduler) Stop() {
	if atomic.CompareAndSwapInt32(&s.status, 0, 1) {
		s.c.Stop()
		s.c = cron.New()
	}
}
