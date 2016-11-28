package logger

import "testing"
import "time"
import "strings"

func TestLog(tx *testing.T) {
	manager.isClose = true
	t, err := time.Parse("2006/01/02 15:04:05", "2016/11/28 16:38:27")
	if err != nil {
		tx.Errorf("test fail, %+v", err)
	}
	manager.Log(LogEvent{Level: "Info", Now: t, Name: "name1", Session: "12345678", Content: "content1", Output: "output1"}, nil)

	// 测试完成打开appender，否则影响其他测试
	manager.isClose = false

	// 写入一个类型不存在的日志，进入记录系统日志的方法
	manager.Log(LogEvent{Level: "test", Now: t, Name: "name1", Session: "12345678", Content: "content1", Output: "output1"}, func(err error) {
		if !strings.EqualFold("不支持的日志类型:test", err.Error()) {
			tx.Errorf("test fail:%v", err)
		}
	})
}

// 单独测试clearUp中的关键代码
func (a *loggerManager) testclearUp() {
	count := a.appenders.RemoveIterCb(func(key string, v interface{}) bool {
		apd := v.(*appenderEntity)
		if time.Now().Sub(apd.last).Seconds() > 5 {
			apd.appender.Close()
			return true
		}
		return false
	})
	if count > 0 {
		sysLoggerInfo(nil, "已移除:", count)
	}
}

func TestClearUp(tx *testing.T) {
	// 保证至少有一个appender
	t, err := time.Parse("2006/01/02 15:04:05", "2016/11/28 16:38:27")
	if err != nil {
		tx.Errorf("test fail, %+v", err)
	}
	manager.Log(LogEvent{Level: "Info", Now: t, Name: "name1", Session: "12345678", Content: "content1", Output: "output1"}, nil)
	count := len(manager.appenders.Keys())

	// 休眠6秒
	time.Sleep(time.Second * 6)

	// 调用方法，判断是否被清理
	manager.testclearUp()
	if len(manager.appenders.Keys()) != 0 {
		tx.Errorf("test fail before count:%d, now:%d", count, len(manager.appenders.Keys()))
	}
}

func TestManagerClose(tx *testing.T) {
	t, err := time.Parse("2006/01/02 15:04:05", "2016/11/28 16:38:27")
	if err != nil {
		tx.Errorf("test fail, %+v", err)
	}
	manager.Log(LogEvent{Level: "Info", Now: t, Name: "name1", Session: "12345678", Content: "content1", Output: "output1"}, nil)
	if len(manager.appenders.Keys()) == 0 {
		tx.Error("test fail")
	}

	manager.Close()
	if len(manager.appenders.Keys()) != 0 {
		tx.Errorf("test fail:manager appenders have:%d", len(manager.appenders.Keys()))
	}
}
