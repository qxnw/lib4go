package logger

import (
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"
)

type TestType struct {
	name string
	age  int
}

func TestNew(t *testing.T) {
	// 日志对象名正确
	log := New("key")
	if len(log.names) != 1 || len(log.sessions) != 1 || len(log.sessions[0]) != 8 {
		t.Error("test fail")
	}

	// 日志对象名为空字符串
	log = New("")
	if len(log.names) != 1 || len(log.sessions) != 1 || len(log.sessions[0]) != 8 {
		t.Error("test fail")
	}

	// 日志对象名为特殊字符
	log = New("!@#!")
	if len(log.names) != 1 || len(log.sessions) != 1 || len(log.sessions[0]) != 8 {
		t.Error("test fail")
	}

	// 多个name
	log = New("key1", "key2")
	if len(log.names) != 2 || len(log.sessions) != 2 || len(log.sessions[0]) != 8 || len(log.sessions[1]) != 8 {
		t.Error("test fail")
	}

	// 判断生成的顺序是否相同
	log = New("key0", "key1", "key2", "key3", "key4", "key5")
	for i, name := range log.names {
		key := fmt.Sprintf("%s%d", "key", i)
		if !strings.EqualFold(name, key) {
			t.Error("test fail")
		}
	}

	// name相同
	log = New("key1", "key1")
	if len(log.names) != 2 || len(log.sessions) != 2 || len(log.sessions[0]) != 8 || len(log.sessions[1]) != 8 {
		t.Error("test fail")
	}

	// 输入空names
	log = New(nil...)
	if len(log.names) != 0 || len(log.sessions) != 0 {
		t.Error("test fail")
	}
}

func TestGet(t *testing.T) {
	// 创建一个日志组件
	// 获取日志组件，判断session id是否为8位
	log := Get("key")
	if len(log.names) != 1 || len(log.sessions) != 1 || len(log.sessions[0]) != 8 || !strings.EqualFold(log.names[0], "key") {
		t.Error("test fail")
	}
	session := log.sessions[0]

	// 判断session id 是否相同
	log = Get("key")
	if len(log.names) != 1 || len(log.sessions) != 1 || len(log.sessions[0]) != 8 || !strings.EqualFold(log.names[0], "key") {
		t.Error("test fail")
	}
	if !strings.EqualFold(log.sessions[0], session) {
		t.Error("test fail")
	}

	// 清空loggers
	loggers.Clear()

	// 获取日志组件，判断session id是否重写创建
	log = Get("key")
	if len(log.names) != 1 || len(log.sessions) != 1 || len(log.sessions[0]) != 8 || !strings.EqualFold(log.names[0], "key") {
		t.Error("test fail")
	}
	if strings.EqualFold(log.sessions[0], session) {
		t.Error("test fail")
	}

	// 日志组件name为空字符串
	log = Get("")
	if len(log.names) != 1 || len(log.sessions) != 1 || len(log.sessions[0]) != 8 || !strings.EqualFold(log.names[0], "") {
		t.Error("test fail")
	}

	// 日志组件name包含特殊字符串
	log = Get("!@#$!%")
	if len(log.names) != 1 || len(log.sessions) != 1 || len(log.sessions[0]) != 8 || !strings.EqualFold(log.names[0], "!@#$!%") {
		t.Error("test fail")
	}

	// 包含多个日志组件名
	log = Get("key1", "key2")
	if len(log.names) != 2 || len(log.sessions) != 2 || len(log.sessions[0]) != 8 || len(log.sessions[1]) != 8 {
		t.Error("test fail")
	}

	// 判断生成的顺序是否相同
	log = Get("key0", "key1", "key2", "key3", "key4", "key5")
	for i, name := range log.names {
		key := fmt.Sprintf("%s%d", "key", i)
		if !strings.EqualFold(name, key) {
			t.Error("test fail")
		}
	}

	// name相同
	log = Get("key1", "key1")
	if len(log.names) != 2 || len(log.sessions) != 2 || len(log.sessions[0]) != 8 || len(log.sessions[1]) != 8 {
		t.Error("test fail")
	}

	// 输入空names
	log = Get(nil...)
	if len(log.names) != 0 || len(log.sessions) != 0 {
		t.Error("test fail")
	}

}

func TestGetSession(t *testing.T) {
	// name, session为正常字符串
	log := GetSession("key", "12345678")
	if len(log.names) != 1 || len(log.sessions) != 1 || !strings.EqualFold(log.names[0], "key") || !strings.EqualFold(log.sessions[0], "12345678") {
		t.Error("test fail")
	}

	// name为空字符串， session为正常字符串
	log = GetSession("", "12345678")
	if len(log.names) != 1 || len(log.sessions) != 1 || !strings.EqualFold(log.names[0], "") || !strings.EqualFold(log.sessions[0], "12345678") {
		t.Error("test fail")
	}

	// name包含特殊字符， session为正常字符串
	log = GetSession("！@#！", "12345678")
	if len(log.names) != 1 || len(log.sessions) != 1 || !strings.EqualFold(log.names[0], "！@#！") || !strings.EqualFold(log.sessions[0], "12345678") {
		t.Error("test fail")
	}

	// name包含特殊字符， session为空字符串
	log = GetSession("key", "")
	if len(log.names) != 1 || len(log.sessions) != 1 || !strings.EqualFold(log.names[0], "key") || !strings.EqualFold(log.sessions[0], "") {
		t.Error("test fail")
	}

	// name包含特殊字符， session包含特殊字符串
	log = GetSession("key", "！@#！")
	if len(log.names) != 1 || len(log.sessions) != 1 || !strings.EqualFold(log.names[0], "key") || !strings.EqualFold(log.sessions[0], "！@#！") {
		t.Error("test fail")
	}

	// name， session包含特殊字符串
	log = GetSession("！@#！", "！@#！")
	if len(log.names) != 1 || len(log.sessions) != 1 || !strings.EqualFold(log.names[0], "！@#！") || !strings.EqualFold(log.sessions[0], "！@#！") {
		t.Error("test fail")
	}
}

func TestGetSessionID(t *testing.T) {
	// 随机生成session id(New)
	log := Get("key1", "key2")
	if len(log.names) != 2 || len(log.sessions) != 2 || len(log.sessions[0]) != 8 || len(log.sessions[1]) != 8 {
		t.Error("test fail")
	}
	if !strings.EqualFold(log.sessions[0], log.GetSessionID()) {
		t.Error("test fail")
	}

	// 手动输入session id(Get)
	log = GetSession("key1", "session1")
	if len(log.names) != 1 || len(log.sessions) != 1 || !strings.EqualFold(log.names[0], "key1") || !strings.EqualFold(log.sessions[0], "session1") {
		t.Error("test fail")
	}
	if !strings.EqualFold(log.sessions[0], log.GetSessionID()) {
		t.Error("test fail")
	}

	// 产生一个空的日志组件
	log = New(nil...)
	if len(log.names) != 0 || len(log.sessions) != 0 {
		t.Error("test fail")
	}
	if !strings.EqualFold(log.GetSessionID(), "") {
		t.Error("test fail")
	}
}

func TestDebug(t *testing.T) {
	// 清空数据统计
	manager.factory = &testLoggerAppenderFactory{}
	ResultClear()

	log := New("key1")

	// 写入字符串
	log.Debug("content1")

	// 写入nil
	log.Debug(nil)

	// 写入int
	log.Debug(1)

	// 写入sliens
	log.Debug(make([]string, 2))

	// 写入数组
	log.Debug([3]int{1, 2, 3})

	// 写入结构体
	log.Debug(TestType{name: "test", age: 11})

	// 日志组件为空
	log = New(nil...)
	log.Debug("hello world")

	time.Sleep(time.Second * 2)
	// 统计数据是否和预期的一致
	if GetResult("debug") != 6 {
		t.Errorf("test fail except : %d, actual : %d", 6, GetResult("debug"))
	}

	Close()
	manager = newLoggerManager()
}

func TestDebugf(t *testing.T) {
	// 清空数据统计
	manager.factory = &testLoggerAppenderFactory{}
	ResultClear()

	log := New("key1")
	// 参数正确
	log.Debugf("%s %s", "hello", "world")

	// format 为空字符串
	log.Debugf("", "hello")

	// format 不包含格式化参数
	log.Debugf("hello", "world")

	// format 格式化参数过多
	log.Debugf("%s %s %s", "hello", "world")

	// 内容为nil
	log.Debugf("hello", nil)

	// 内容和格式化参数类型不匹配
	log.Debugf("%s %d", "hello", "world")

	// 内容为结构体
	log.Debugf("%+v", TestType{name: "test", age: 11})

	// 日志组件为空
	log = New(nil...)
	log.Debugf("%s %s", "hello", "world")

	time.Sleep(time.Second * 2)
	// 统计数据是否和预期的一致
	if GetResult("debug") != 7 {
		t.Errorf("test fail except : %d, actual : %d", 7, GetResult("debug"))
	}

	Close()
	manager = newLoggerManager()
}

func TestInfo(t *testing.T) {
	// 清空数据统计
	manager.factory = &testLoggerAppenderFactory{}
	ResultClear()

	log := New("key1")

	// 写入字符串
	log.Info("content1")
	// // 每秒钟写入文件一次
	// time.Sleep(time.Second * 2)

	// 写入nil
	log.Info(nil)

	// 写入int
	log.Info(1)

	// 写入sliens
	log.Info(make([]string, 2))

	// 写入数组
	log.Info([3]int{1, 2, 3})

	// 写入结构体
	log.Info(TestType{name: "test", age: 11})

	// 日志组件为空
	log = New(nil...)
	log.Info("hello world")

	time.Sleep(time.Second * 2)
	// 统计数据是否和预期的一致
	if GetResult("info") != 6 {
		t.Errorf("test fail except : %d, actual : %d", 6, GetResult("info"))
	}

	Close()
	manager = newLoggerManager()
}

func TestInfof(t *testing.T) {
	// 清空数据统计
	manager.factory = &testLoggerAppenderFactory{}
	ResultClear()

	log := New("key1")
	// 参数正确
	log.Infof("%s %s", "hello", "world")

	// format 为空字符串
	log.Infof("", "hello")

	// format 不包含格式化参数
	log.Infof("hello", "world")

	// format 格式化参数过多
	log.Infof("%s %s %s", "hello", "world")

	// 内容为nil
	log.Infof("hello", nil)

	// 内容和格式化参数类型不匹配
	log.Infof("%s %d", "hello", "world")

	// 内容为结构体
	log.Infof("%+v", TestType{name: "test", age: 11})

	// 日志组件为空
	log = New(nil...)
	log.Infof("%s %s", "hello", "world")

	time.Sleep(time.Second * 2)
	// 统计数据是否和预期的一致
	if GetResult("info") != 7 {
		t.Errorf("test fail except : %d, actual : %d", 7, GetResult("info"))
	}

	Close()
	manager = newLoggerManager()
}

func TestError(t *testing.T) {
	// 清空数据统计
	manager.factory = &testLoggerAppenderFactory{}
	ResultClear()

	log := New("key1")

	// 写入字符串
	log.Error("content1")

	// 写入nil
	log.Error(nil)

	// 写入int
	log.Error(1)

	// 写入sliens
	log.Error(make([]string, 2))

	// 写入数组
	log.Error([3]int{1, 2, 3})

	// 写入结构体
	log.Error(TestType{name: "test", age: 11})

	// 日志组件为空
	log = New(nil...)
	log.Error("hello world")

	time.Sleep(time.Second * 2)
	// 统计数据是否和预期的一致
	if GetResult("error") != 6 {
		t.Errorf("test fail except : %d, actual : %d", 6, GetResult("error"))
	}

	Close()
	manager = newLoggerManager()
}

func TestErrorf(t *testing.T) {
	// 清空数据统计
	manager.factory = &testLoggerAppenderFactory{}
	ResultClear()

	log := New("key1")
	// 参数正确
	log.Errorf("%s %s", "hello", "world")

	// format 为空字符串
	log.Errorf("", "hello")

	// format 不包含格式化参数
	log.Errorf("hello", "world")

	// format 格式化参数过多
	log.Errorf("%s %s %s", "hello", "world")

	// 内容为nil
	log.Errorf("hello", nil)

	// 内容为结构体
	log.Errorf("%+v", TestType{name: "test", age: 11})

	// 内容和格式化参数类型不匹配
	log.Errorf("%s %d", "hello", "world")

	// 日志组件为空
	log = New(nil...)
	log.Errorf("%s %s", "hello", "world")

	time.Sleep(time.Second * 2)
	// 统计数据是否和预期的一致
	if GetResult("error") != 7 {
		t.Errorf("test fail except : %d, actual : %d", 7, GetResult("error"))
	}

	Close()
	manager = newLoggerManager()
}

func doWrite(ch chan int, lk *sync.WaitGroup) {
	log := New("abc")
START:
	for {
		select {
		case v, ok := <-ch:
			if ok {
				log.Debug(v)
				log.Info(v)
				log.Error(v)
			} else {
				break START
			}
		}
	}
	lk.Done()
}

func TestALL(t *testing.T) {
	manager.factory = &testLoggerAppenderFactory{}
	// 清空结果
	ResultClear()
	totalCount := 10000 * 5
	ch := make(chan int, totalCount)
	lk := sync.WaitGroup{}

	for i := 0; i < 100; i++ {
		lk.Add(1)
		go doWrite(ch, &lk)
	}

	for i := 0; i < totalCount; i++ {
		ch <- i
	}

	close(ch)
	lk.Wait()

	// time.Sleep(time.Second * 2)

	Close()

	for i := 0; i < len(ACCOUNT); i++ {
		fmt.Println(ACCOUNT[i].name, " ", ACCOUNT[i].count)
		if strings.EqualFold(ACCOUNT[i].name, "debug") {
			if ACCOUNT[i].count != totalCount {
				t.Errorf("test fail, actual : %d", ACCOUNT[i].count)
			}
		}
		if strings.EqualFold(ACCOUNT[i].name, "info") {
			if ACCOUNT[i].count != totalCount {
				t.Errorf("test fail, actual : %d", ACCOUNT[i].count)
			}
		}
		// 测试不执行fatal日志记录
		if strings.EqualFold(ACCOUNT[i].name, "fatal") {
			if ACCOUNT[i].count != 0 {
				t.Errorf("test fail, actual : %d", ACCOUNT[i].count)
			}
		}
		if strings.EqualFold(ACCOUNT[i].name, "error") {
			if ACCOUNT[i].count != totalCount {
				t.Errorf("test fail, actual : %d", ACCOUNT[i].count)
			}
		}
	}

	manager = newLoggerManager()
}
