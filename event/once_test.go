package event

import "testing"
import "fmt"
import "time"

// TestNewOnce 测试构建一次执行锁
func TestNewOnce(t *testing.T) {
	// 构建一个一次执行锁
	_, err := NewOnce(1)
	if err != nil {
		t.Errorf("test fail %v", err)
	}

	_, err = NewOnce(-1)
	if err == nil {
		t.Error("test fail")
	}

	_, err = NewOnce(0)
	if err == nil {
		t.Error("test fail")
	}
}

// TestWait 测试等待所有的任务执行
func TestWait(t *testing.T) {
	// 构建一个一次性执行锁
	n := 3
	sn, err := NewOnce(n)
	if err != nil {
		t.Errorf("test fail %v", err)
	}

	// 等待完成标志
	isFinish := false

	go func() {
		sn.Wait()
		isFinish = true
	}()

	for i := 0; i < n; i++ {
		sn.Done(fmt.Sprintf("%d", i))
		time.Sleep(time.Second)
	}

	if !isFinish {
		t.Errorf("test fail")
	}
}

// TestAddStep 测试添加步骤
func TestAddStep(t *testing.T) {
	// 构建一个一次性执行锁
	n := 3
	sn, err := NewOnce(n)
	if err != nil {
		t.Errorf("test fail %v", err)
	}

	// 等待完成标志
	isFinish := false

	go func() {
		sn.Wait()
		isFinish = true
	}()

	for i := 0; i < n; i++ {
		if i == n-1 {
			sn.AddStep(1)
		}
		sn.Done(fmt.Sprintf("%d", i))
		time.Sleep(time.Second)
	}

	if isFinish {
		t.Errorf("test fail")
	}

	sn.Done(time.Now().String())

	if !isFinish {
		t.Error("test fail")
	}
}
