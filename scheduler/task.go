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
	// defer func() {
	// 	if err := recover(); err != nil {
	// 		fmt.Println(err)
	// 	}
	// }()
	// go j.fun(j.obj)

	/*change by champly 2016年11月22日12:00:41*/
	// 捕获异常
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()
		j.fun(j.obj)
	}()
	/*end*/
}
