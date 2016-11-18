package sysinfo

import (
	"fmt"
	"testing"
)

func TestGetSysinfo(t *testing.T) {
	fmt.Println(GetAPPMemory())
	fmt.Println(GetMemory())
	fmt.Println(GetCPU())
	fmt.Println(GetDisk())
}
