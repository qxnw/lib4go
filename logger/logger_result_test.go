package logger

import "strings"
import "sync"

type Account struct {
	name  string
	count int
}

var mutex sync.Mutex

var ACCOUNT []*Account

// SetResult 存放测试结果
func SetResult(name string, n int) {
	for i := 0; i < len(ACCOUNT); i++ {
		if strings.EqualFold(ACCOUNT[i].name, name) {
			mutex.Lock()
			ACCOUNT[i].count = ACCOUNT[i].count + n
			mutex.Unlock()
			return
		}
	}

	mutex.Lock()
	account := &Account{name: name, count: n}
	ACCOUNT = append(ACCOUNT, account)
	mutex.Unlock()
}

// GetResult 获取测试结果
func GetResult(name string) int {
	for i := 0; i < len(ACCOUNT); i++ {
		if strings.EqualFold(ACCOUNT[i].name, name) {
			return ACCOUNT[i].count
		}
	}
	return 0
}
