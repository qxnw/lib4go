package scheduler

import (
	"fmt"
	"testing"
	"time"
)

func TestNewScheduler(t *testing.T) {
	// 正常流程
	scheduler := NewScheduler()
	trigger := "0/1 * * * * ?"
	notifyChan := make(chan int)
	taskDetail := NewTask(notifyChan, func(v interface{}) {
		fmt.Println("start func")
		nchan := v.(chan int)
		nchan <- 1
		fmt.Println(<-nchan)
	})
	scheduler.AddTask(trigger, taskDetail)
	scheduler.Start()
	// 休眠等待调用
	time.Sleep(time.Second * 2)
	scheduler.Stop()

	AddTask(trigger, taskDetail)
	Start()
	if Count() != 1 {
		t.Error("test fail Count:", Count())
	}
	Stop()
	if Count() != 0 {
		t.Error("test fail Count:", Count())
	}

	// trigger错误
	scheduler = NewScheduler()
	trigger = "1231321"
	notifyChan = make(chan int)
	taskDetail = NewTask(notifyChan, func(v interface{}) {
		fmt.Println("start func")
		nchan := v.(chan int)
		nchan <- 1
		fmt.Println(<-nchan)
	})
	scheduler.AddTask(trigger, taskDetail)
	scheduler.Start()
	time.Sleep(time.Second * 2)
	scheduler.Stop()

	AddTask(trigger, taskDetail)
	Start()
	time.Sleep(time.Second * 2)
	if Count() != 0 {
		t.Error("test fail Count:", Count())
	}
	time.Sleep(time.Second * 2)
	Stop()
	if Count() != 0 {
		t.Error("test fail Count:", Count())
	}

	// 传入空的Job
	scheduler = NewScheduler()
	trigger = "0/1 * * * * ?"
	var obj interface{}
	taskDetail = NewTask(obj, func(v interface{}) {
		fmt.Println("start func")
	})
	scheduler.AddTask(trigger, taskDetail)
	scheduler.Start()
	time.Sleep(time.Second * 2)
	scheduler.Stop()

	AddTask(trigger, taskDetail)
	Start()
	time.Sleep(time.Second * 2)
	if Count() != 1 {
		t.Error("test fail Count:", Count())
	}
	Stop()
	if Count() != 0 {
		t.Error("test fail Count:", Count())
	}

	// 异常调用
	taskDetail = NewTask(obj, func(v interface{}) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()
		panic("this is a panic")
	})

	taskDetail.Run()
	time.Sleep(time.Second * 5)
}
