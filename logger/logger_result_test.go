package logger

import "strings"

type Account struct {
	name  string
	count int
}

var ACCOUNT []*Account

// SetResult 存放测试结果
func SetResult(name string) {
	for i := 0; i < len(ACCOUNT); i++ {
		if strings.EqualFold(ACCOUNT[i].name, name) {
			ACCOUNT[i].count++
			return
		}
	}

	account := &Account{name: name, count: 1}
	ACCOUNT = append(ACCOUNT, account)
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
