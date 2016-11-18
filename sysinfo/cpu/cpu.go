package cpu

import (
	"encoding/json"
	"fmt"

	"github.com/shirou/gopsutil/cpu"
)

// Useage Total总量，Idle空闲，Used使用率，Collercter总量，使用量
type Useage struct {
	Total     string `json:"total"`
	Idle      string `json:"idle"`
	Used      string `json:"used"`
	Collecter []string
}

// GetInfo 获取当前系统CPU数据
func GetInfo() []map[string]interface{} {
	v, _ := cpu.Times(true)
	buffer, _ := json.Marshal(&v)
	var data []map[string]interface{}
	json.Unmarshal(buffer, &data)
	return data
}

// GetAvaliabeInfo 获取当前系统CPU使用的情况数据
func GetAvaliabeInfo() (useage Useage) {
	cpus, _ := cpu.Times(true)
	useage = Useage{}
	var total, idle float64
	for _, value := range cpus {
		total += value.Total()
		idle += value.Idle
	}
	useage.Total = fmt.Sprintf("%.2f", total)
	useage.Idle = fmt.Sprintf("%.2f", idle)
	useage.Used = fmt.Sprintf("%.2f", (total-idle)/total)
	useage.Collecter = []string{useage.Total, useage.Used}
	return
}
