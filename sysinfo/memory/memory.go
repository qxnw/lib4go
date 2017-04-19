package memory

import "github.com/shirou/gopsutil/mem"

// Useage Total总量，Idle空闲，Used使用率，Collercter总量，使用量
type Useage struct {
	Total       uint64  `json:"total"`
	Idle        uint64  `json:"idle"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"used"`
}

// GetInfo 获取当前系统内存使用数据
func GetInfo() (useage Useage) {
	vm, _ := mem.VirtualMemory()
	useage.Total = vm.Total
	useage.Idle = vm.Free
	useage.Used = vm.Used
	useage.UsedPercent = vm.UsedPercent
	return
}
