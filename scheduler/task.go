package scheduler

import "fmt"

type TaskDetail struct {
	obj interface{}
	fun func(obj interface{})
}

func NewTask(obj interface{}, fun func(obj interface{})) *TaskDetail {
	return &TaskDetail{obj: obj, fun: fun}
}

func (j *TaskDetail) Run() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	go j.fun(j.obj)
}
