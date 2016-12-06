package zk

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

// TestSepcialSituation 测试特殊情况下zookeeper对应的错误处理
func TestSepcialSituation(t *testing.T) {
	// 连接到zookeeper服务器
	address := "192.168.0.159:2181"
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

	fmt.Println("断开连接")

	time.Sleep(time.Second * 20)

	fmt.Println("恢复连接")

	time.Sleep(time.Second * 30)

	client.Disconnect()
}

// TestBadNetworkSituation 测试网络从异常中恢复之后获取节点的值是否异常
func TestBadNetworkSituation(t *testing.T) {
	// 连接到zookeeper服务器
	path := "/zk_test/123"
	address := "192.168.0.159:2181"
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

	fmt.Println("断开连接")

	time.Sleep(time.Second * 20)

	// 网络断开
	{
		// 获取节点的值
		_, err = client.GetValue(path)
		if !strings.EqualFold(err.Error(), "未连接到zk服务器") {
			t.Error("test fail")
		}

		// 获取节点的子节点
		_, err = client.GetChildren(path)
		if !strings.EqualFold(err.Error(), "未连接到zk服务器") {
			t.Error("test fail")
		}

	}

	fmt.Println("恢复连接")

	time.Sleep(time.Second * 30)

	// 网络恢复
	{
		// 获取节点的值
		value, err := client.GetValue(path)
		if err != nil {
			t.Errorf("test fail %v", err)
		}
		if !strings.EqualFold(value, "test") {
			t.Errorf("test fail actual : %s, except:%s", value, "test")
		}

		// 获取节点的子节点
		paths, err := client.GetChildren(path)
		if err != nil {
			t.Errorf("test fail %v", err)
		}
		if len(paths) != 0 {
			t.Error("test fail")
		}
	}

	client.Disconnect()
}

// TestBadNetworkBindWatchValue 测试网络断开的情况下获取节点的值
func TestBadNetworkBindWatchValue(t *testing.T) {
	// 连接到zookeeper服务器
	path := "/zk_test/123"
	address := "192.168.0.159:2181"
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

	// 监控节点的变化
	data := make(chan string, 1)
	go func() {
		err = client.BindWatchValue(path, data)
		fmt.Println(err)
	}()

	fmt.Println("修改/zk_test/123节点的值")
	fmt.Println(<-data)

	fmt.Println("断开连接")
	time.Sleep(time.Second * 10)

	fmt.Println("恢复连接")
	fmt.Println(<-data)

	client.Disconnect()
}

// TestNoNetworkUnBindWatchValue 测试在没有网络的时候取消绑定
func TestNoNetworkUnBindWatchValue(t *testing.T) {
	// 连接到zookeeper服务器
	path := "/zk_test/123"
	address := "192.168.0.159:2181"
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

	// 监控节点的变化
	data := make(chan string, 1)
	go func() {
		err = client.BindWatchValue(path, data)
		fmt.Println(err)
	}()

	fmt.Println("修改/zk_test/123节点的值")
	fmt.Println(<-data)

	fmt.Println("断开连接")
	time.Sleep(time.Second * 10)

	fmt.Println("取消绑定")
	client.UnbindWatchValue(path)

	fmt.Println("恢复连接")
	// fmt.Println(<-data)

	time.Sleep(time.Second * 20)

	client.Disconnect()
}

// TestBadNetWorkBindWatchChildren 测试网络异常的情况下监控节点下面的子节点
func TestBadNetWorkBindWatchChildren(t *testing.T) {
	// 连接到zookeeper服务器
	path := "/zk_test"
	address := "192.168.0.159:2181"
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

	// 监控节点的变化
	data := make(chan []string, 1)
	go func() {
		err = client.BindWatchChildren(path, data)
		fmt.Println(err)
	}()

	fmt.Println("修改/zk_test节点的子节点")
	fmt.Println(<-data)

	fmt.Println("断开连接")
	time.Sleep(time.Second * 10)

	fmt.Println("恢复连接")
	fmt.Println("修改/zk_test的子节点")
	fmt.Println(<-data)

	client.Disconnect()
}

// TestNoNetworkUnBindWatchChildren 测试网络异常情况下取消绑定
func TestNoNetworkUnBindWatchChildren(t *testing.T) {
	// 连接到zookeeper服务器
	path := "/zk_test"
	address := "192.168.0.159:2181"
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

	// 监控节点的变化
	data := make(chan []string, 1)
	go func() {
		err = client.BindWatchChildren(path, data)
		fmt.Println(err)
	}()

	fmt.Println("修改/zk_test节点的子节点")
	fmt.Println(<-data)

	fmt.Println("断开连接")
	time.Sleep(time.Second * 10)
	fmt.Println("取消绑定")
	client.UnbindWatchChildren(path)

	fmt.Println("恢复连接")
	// fmt.Println("修改/zk_test的子节点")
	// fmt.Println(<-data)

	client.Disconnect()
}
