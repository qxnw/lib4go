package disk

import (
	"fmt"
	"runtime"

	"github.com/shirou/gopsutil/disk"
)

// Useage Total总量，Idle空闲，Used使用率，Collercter总量，使用量
type Useage struct {
	Total     string `json:"total"`
	Idle      string `json:"idle"`
	Used      string `json:"used"`
	Collecter []string
}

// GetInfo 获取磁盘使用信息
func GetInfo() (useage Useage) {
	dir := "/"
	if runtime.GOOS == "windows" {
		dir = "c:"
	}
	sm, _ := disk.Usage(dir)
	useage.Total = fmt.Sprintf("%d", sm.Total)
	useage.Idle = fmt.Sprintf("%d", sm.Total-sm.Used)
	useage.Used = fmt.Sprintf("%.2f", sm.UsedPercent)
	useage.Collecter = []string{useage.Total, useage.Used}
	return
}
