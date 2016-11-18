package disk

import (
	"encoding/json"
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

// GetInfo 获取磁盘信息
func GetInfo() (data []map[string]interface{}) {
	data = make([]map[string]interface{}, 0, 6)
	defer func() {
		if er := recover(); er != nil {
			fmt.Println("get DISK error", er)
		}
	}()
	var stats []*disk.UsageStat
	if runtime.GOOS == "windows" {
		v, _ := disk.Partitions(true)
		for _, p := range v {
			s, _ := disk.Usage(p.Device)
			stats = append(stats, s)
		}
	} else {
		s, _ := disk.Usage("/")
		stats = append(stats, s)
	}

	buffer, _ := json.Marshal(&stats)
	json.Unmarshal(buffer, &data)
	return
}

// GetAvaliabeInfo 获取磁盘使用信息
func GetAvaliabeInfo() (useage Useage) {
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
