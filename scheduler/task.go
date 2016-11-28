package scheduler

import "fmt"

//Task 任务信息
type Task struct {
	input []interface{}
	fun   func(input ...interface{})
}

//NewTask 创建任务
func NewTask(fun func(input ...interface{}), input ...interface{}) *Task {
	return &Task{input: input, fun: fun}
}

//Run 执行任务
func (j *Task) Run() {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()
		j.fun(j.input...)
	}()
}
