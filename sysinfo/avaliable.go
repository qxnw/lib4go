package sysinfo

import (
	"fmt"

	"runtime"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

type Useage struct {
	Total     string `json:"total"`
	Idle      string `json:"idle"`
	Used      string `json:"used"`
	Collecter []string
}

func GetAvaliabeCPU() (useage Useage) {
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

func GetAvaliabeMem() (useage Useage) {
	vm, _ := mem.VirtualMemory()
	useage.Total = fmt.Sprintf("%d", vm.Total)
	useage.Idle = fmt.Sprintf("%d", vm.Total-vm.Used)
	useage.Used = fmt.Sprintf("%.2f", vm.UsedPercent)
	useage.Collecter = []string{useage.Total, useage.Used}
	return
}

func GetAvaliabeDisk() (useage Useage) {
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
