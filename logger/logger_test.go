package logger

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	manager.factory = &testLoggerAppenderFactory{}

	log := New("logger")

	//session := log.GetSessionID()
	//if len(session) != 8 {
	//	t.Error("test fail")
	//}

	//log = Get("newlogger")
	//log = GetSession("newlogger", log.GetSessionID())

	log.Debug("hello world")
	/*log.Debugf("%s %s", "hello", "world")
	log.Info("hello world")
	log.Infof("%s %s", "hello", "world")
	*/
	time.Sleep(time.Second * 5)
	//log.Debug("timeout")

	Close()

	//for i := 0; i < len(ACCOUNT); i++ {
	//fmt.Println(ACCOUNT[i].name, " ", ACCOUNT[i].count)
	//}

	// n := 100
	// for i := 0; i < n; i++ {
	// 	go func(i int) {
	// 		l := New("test", "test1")
	// 		for j := 0; j < 5; j++ {
	// 			l.Debug("当前数量:", i)
	// 			time.Sleep(time.Second)
	// 		}
	// 	}(i)
	// }

	// time.Sleep(time.Second * 20)
	// // Close()

	// for i := 0; i < n; i++ {
	// 	go func(i int) {
	// 		l := Get("test", "test1")
	// 		for j := 0; j < 5; j++ {
	// 			l.Debug("当前数量:", i)
	// 			time.Sleep(time.Second)
	// 		}
	// 	}(i)
	// }

	// time.Sleep(time.Second * 20)
	// Close()
}
