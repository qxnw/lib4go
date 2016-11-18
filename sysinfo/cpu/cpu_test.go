package cpu

import "testing"

func Test(t *testing.T) {
	info := GetInfo()
	if info == nil {
		t.Error("test fail")
	}
	// for _, cpu := range info {
	// 	for key, value := range cpu {
	// 		fmt.Println(key, " ", value)
	// 	}
	// }

	avaliable := GetAvaliabeInfo()
	if avaliable.Total == "" || avaliable.Idle == "" || avaliable.Used == "" {
		t.Error("test fail")
	}
}
