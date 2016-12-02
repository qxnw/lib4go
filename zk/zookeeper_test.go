package zk

import (
	"strings"
	"testing"
	"time"
)

var address = "192.168.0.159:2181"

// TestNew 测试创建zookeeper连接
func TestNew(t *testing.T) {
	timeout := time.Second * 1
	loggerName := "zookeeper"

	// 正确的ip地址
	servers := []string{address}
	_, err := New(servers, timeout, loggerName)
	if err != nil {
		t.Errorf("test fail %v", err)
	}
}

// TestConnect 测试连接到zookeeper服务器
func TestConnect(t *testing.T) {
	timeout := time.Second * 1
	loggerName := "zookeeper"

	// 正确的ip地址
	servers := []string{address}
	client, err := New(servers, timeout, loggerName)
	if err != nil {
		t.Errorf("test fail %v", err)
	}
	err = client.Connect()
	if err != nil {
		t.Errorf("test fail %v", err)
	}

	time.Sleep(time.Second * 1)

	if !client.isConnect {
		t.Error("test fail")
	}
	client.Disconnect()

	// ip地址不对
	servers = []string{"192.168.0.165:2181"}
	client, err = New(servers, timeout, loggerName)
	if err != nil {
		t.Errorf("test fail %v", err)
	}
	err = client.Connect()
	if err != nil {
		t.Errorf("test fail %v", err)
	}

	time.Sleep(time.Second * 1)

	if client.isConnect {
		t.Error("test fail")
	}

	// 端口不对
	servers = []string{"192.168.0.166:12181"}
	client, err = New(servers, timeout, loggerName)
	if err != nil {
		t.Errorf("test fail %v", err)
	}
	err = client.Connect()
	if err != nil {
		t.Errorf("test fail %v", err)
	}

	time.Sleep(time.Second * 1)

	if client.isConnect {
		t.Error("test fail")
	}

	// ip地址格式不对
	servers = []string{"asdfaqe"}
	client, err = New(servers, timeout, loggerName)
	if err != nil {
		t.Errorf("test fail %v", err)
	}
	err = client.Connect()
	if err == nil {
		t.Errorf("test fail %v", err)
	}
}

// TestReconnect 测试重新连接服务器
func TestReconnect(t *testing.T) {
	timeout := time.Second * 1
	loggerName := "zookeeper"
	servers := []string{address}

	// 没有连接过，直接重连
	client, err := New(servers, timeout, loggerName)
	if err != nil {
		t.Errorf("test fail %v", err)
	}
	err = client.Reconnect()
	if err != nil {
		t.Errorf("test fail %v", err)
	}

	time.Sleep(time.Second * 1)

	if !client.isConnect {
		t.Error("test fail")
	}
	client.Disconnect()

	// 连接之后，关闭连接，重新连接
	client, err = New(servers, timeout, loggerName)
	if err != nil {
		t.Errorf("test fail %v", err)
	}
	err = client.Connect()
	if err != nil {
		t.Errorf("test fail %v", err)
	}

	time.Sleep(time.Second * 1)

	if !client.isConnect {
		t.Error("test fail")
	}
	// 关闭连接
	client.Disconnect()
	if client.isConnect {
		t.Error("test fail")
	}

	// 重新连接
	err = client.Reconnect()
	if err != nil {
		t.Errorf("test fail %v", err)
	}

	time.Sleep(time.Second * 1)

	if !client.isConnect {
		t.Error("test fail")
	}
}

// TestDisconnect 测试关闭zookeeper连接
func TestDisconnect(t *testing.T) {
	timeout := time.Second * 1
	loggerName := "zookeeper"
	servers := []string{address}

	// 只有一个user的时候
	{
		client, err := New(servers, timeout, loggerName)
		if err != nil {
			t.Errorf("test fail %v", err)
		}
		err = client.Connect()
		if err != nil {
			t.Errorf("test fail %v", err)
		}

		time.Sleep(time.Second * 1)

		if !client.isConnect {
			t.Error("test fail")
		}
		client.Disconnect()

		// 判断结果
		if client.isConnect {
			t.Error("test fail")
		}
	}

	// 没有user的时候
	{
		client, err := New(servers, timeout, loggerName)
		if err != nil {
			t.Errorf("test fail %v", err)
		}
		err = client.Connect()
		if err != nil {
			t.Errorf("test fail %v", err)
		}

		time.Sleep(time.Second * 1)

		if !client.isConnect {
			t.Error("test fail")
		}
		client.Disconnect()
		// 多次关闭
		client.Disconnect()

		// 判断结果
		if client.isConnect {
			t.Error("test fail")
		}
	}

	// 有多个user的时候
	{
		client, err := New(servers, timeout, loggerName)
		if err != nil {
			t.Errorf("test fail %v", err)
		}
		err = client.Connect()
		if err != nil {
			t.Errorf("test fail %v", err)
		}
		// 连接两次
		err = client.Connect()
		if err != nil {
			t.Errorf("test fail %v", err)
		}

		time.Sleep(time.Second * 1)

		if !client.isConnect {
			t.Error("test fail")
		}
		client.Disconnect()

		// 判断结果
		if !client.isConnect {
			t.Error("test fail")
		}
	}

	// 还没有连接过的时候
	{
		client, err := New(servers, timeout, loggerName)
		if err != nil {
			t.Errorf("test fail %v", err)
		}

		client.Disconnect()

		// 判断结果
		if client.isConnect {
			t.Error("test fail")
		}
	}
}

// TestExistsAny 测试是否有一个路径存在
func TestExistsAny(t *testing.T) {
	// zk_test节点存在
	paths := "/zk_test"
	timeout := time.Second * 1
	loggerName := "zookeeper"
	servers := []string{address}
	client, err := New(servers, timeout, loggerName)
	if err != nil {
		t.Errorf("test fail %v", err)
	}
	err = client.Connect()
	if err != nil {
		t.Errorf("test fail %v", err)
	}

	time.Sleep(time.Second * 1)

	if !client.isConnect {
		t.Error("test fail")
	}

	// 有一个节点存在
	b, actual, err := client.ExistsAny(paths, "/home", "/err_path")
	if err != nil {
		t.Errorf("test fail %v", err)
	}
	if !b {
		t.Error("test fail")
	}
	if !strings.EqualFold(actual, paths) {
		t.Errorf("test fail actual : %s, except : %s", actual, paths)
	}

	// 没有存在的节点
	b, actual, err = client.ExistsAny("/home", "/err_path")
	if err != nil {
		t.Errorf("test fail %v", err)
	}
	if b {
		t.Error("test fail")
	}

	// 路径包含特殊字符
	b, actual, err = client.ExistsAny("/home", "！@#！@！@")
	if err != nil {
		t.Errorf("test fail %v", err)
	}
	if b {
		t.Error("test fail")
	}

	client.Disconnect()

	// client 没有连接到zookeeper服务器
	client, err = New(servers, timeout, loggerName)
	if err != nil {
		t.Errorf("test fail %v", err)
	}
	b, actual, err = client.ExistsAny("/home", "！@#！@！@")
	if err == nil {
		t.Errorf("test fail")
	}
	if b {
		t.Error("test fail")
	}
}

// TestExists 测试判断一个节点是否存在
func TestExists(t *testing.T) {
	// zk_test节点存在
	path := "/zk_test"
	timeout := time.Second * 1
	loggerName := "zookeeper"
	servers := []string{address}
	client, err := New(servers, timeout, loggerName)
	if err != nil {
		t.Errorf("test fail %v", err)
	}
	err = client.Connect()
	if err != nil {
		t.Errorf("test fail %v", err)
	}

	time.Sleep(time.Second * 1)

	if !client.isConnect {
		t.Error("test fail")
	}

	// 节点存在
	b, err := client.Exists(path)
	if err != nil {
		t.Errorf("test fail %v", err)
	}
	if !b {
		t.Error("test fail")
	}

	// 节点不存在
	b, err = client.Exists("/home")
	if err != nil {
		t.Errorf("test fail %v", err)
	}
	if b {
		t.Error("test fail")
	}

	// 路径包含特殊字符
	b, err = client.Exists("！@#！@！@")
	if err != nil {
		t.Errorf("test fail %v", err)
	}
	if b {
		t.Error("test fail")
	}

	client.Disconnect()

	// client 没有连接到zookeeper服务器
	client, err = New(servers, timeout, loggerName)
	if err != nil {
		t.Errorf("test fail %v", err)
	}

	b, err = client.Exists(path)
	if err == nil {
		t.Errorf("test fail")
	}
	if b {
		t.Error("test fail")
	}
}

// TestDelete 测试删除一个节点
func TestDelete(t *testing.T) {
	// zk_test节点存在
	path := "/zk_test/123"
	timeout := time.Second * 1
	loggerName := "zookeeper"
	servers := []string{address}
	client, err := New(servers, timeout, loggerName)
	if err != nil {
		t.Errorf("test fail %v", err)
	}
	err = client.Connect()
	if err != nil {
		t.Errorf("test fail %v", err)
	}

	time.Sleep(time.Second * 1)

	if !client.isConnect {
		t.Error("test fail")
	}

	// 删除一个不存在的节点
	err = client.Delete("/err_path/err_node")
	if err == nil {
		t.Error("test fail")
	}

	// 判断节点是否存在
	b, err := client.Exists(path)
	if err != nil {
		t.Errorf("test fail %v", err)
	}
	if !b {
		t.Error("节点不存在")
	} else {
		err = client.Delete(path)

		if err != nil {
			t.Errorf("test fail %v", err)
		}
	}
}

// TestGetPaths 测试获取当前路径下的所有子路径
func TestGetPaths(t *testing.T) {
	// 节点存在

	// 节点不存在

	// 没有子节点
}
