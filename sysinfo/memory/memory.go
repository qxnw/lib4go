package memory

import (
	"fmt"

	"github.com/shirou/gopsutil/mem"
)

// Useage Total总量，Idle空闲，Used使用率，Collercter总量，使用量
type Useage struct {
	Total     string `json:"total"`
	Idle      string `json:"idle"`
	Used      string `json:"used"`
	Collecter []string
}

// GetInfo 获取当前系统内存使用数据
func GetInfo() (useage Useage) {
	vm, _ := mem.VirtualMemory()
	useage.Total = fmt.Sprintf("%d", vm.Total)
	useage.Idle = fmt.Sprintf("%d", vm.Total-vm.Used)
	useage.Used = fmt.Sprintf("%.2f", vm.UsedPercent)
	useage.Collecter = []string{useage.Total, useage.Used}
	return
}
