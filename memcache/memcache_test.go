package memcache

import "testing"

func TestMemcache(tx *testing.T) {
	_, err := New("192.168.0.166:11212")
	if err != nil {
		tx.Error(err)
		return
	}
	/*key := "test"
	value := "value"
	if e := mem.Set(key, value, 10); e != nil {
		t.Errorf("存数据失败：%s", value)
	}

	v := mem.Get(key)
	if !strings.EqualFold(v, value) {
		t.Errorf("取数据失败：%s", v)
	}

	err = mem.Delay(key, 10)
	if err != nil {
		t.Errorf("延长存储时间失败:%v", err)
	}

	time.Sleep(15 * time.Second)

	v = mem.Get(key)
	if !strings.EqualFold(v, value) {
		t.Error("延长存储时间失败")
	}

	mem.Delete(key)
	v = mem.Get(key)
	if !strings.EqualFold(v, "") {
		t.Errorf("获取到了脏数:%s", v)
	}*/

}
