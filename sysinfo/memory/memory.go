package memory

import (
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/shirou/gopsutil/mem"
)

// Useage Total总量，Idle空闲，Used使用率，Collercter总量，使用量
type Useage struct {
	Total     string `json:"total"`
	Idle      string `json:"idle"`
	Used      string `json:"used"`
	Collecter []string
}

// GetAPPInfo 获取App内存信息
func GetAPPInfo() uint64 {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	return mem.Alloc >> 20
}

// GetInfo 获取当前系统的内存信息
func GetInfo() []map[string]interface{} {
	v, _ := mem.VirtualMemory()
	data := make(map[string]interface{})
	buffer, _ := json.Marshal(&v)
	json.Unmarshal(buffer, &data)
	var result []map[string]interface{}
	result = append(result, data)
	return result
}

// GetAvaliabeInfo 获取当前系统内存使用数据
func GetAvaliabeInfo() (useage Useage) {
	vm, _ := mem.VirtualMemory()
	useage.Total = fmt.Sprintf("%d", vm.Total)
	useage.Idle = fmt.Sprintf("%d", vm.Total-vm.Used)
	useage.Used = fmt.Sprintf("%.2f", vm.UsedPercent)
	useage.Collecter = []string{useage.Total, useage.Used}
	return
}
