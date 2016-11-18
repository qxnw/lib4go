package sysinfo

import (
	"fmt"
	"testing"
)

func TestAvaliabe(t *testing.T) {
	fmt.Println(GetAvaliabeCPU())
	fmt.Println(GetAvaliabeMem())
	fmt.Println(GetAvaliabeDisk())

	// data, _ := json.Marshal(GetAvaliabeCPU())
	// fmt.Println(string(data))
}
