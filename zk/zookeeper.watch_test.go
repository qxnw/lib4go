package zk

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

/*
   192.168.0.159:2181
   192.168.0.154:2181
   做的集群，159为主
*/

var (
	masterAddress = "192.168.0.159:2181"
	followAddress = "192.168.0.154:2181"
)

// TestBindWatchValue 测试监控一个节点的值是否发送变化
func TestBindWatchValue(t *testing.T) {
	// master client
	timeout := time.Second * 1
	loggerName := "zookeeper"
	servers := []string{masterAddress}
	masterClient, err := New(servers, timeout, loggerName)
	if err != nil {
		t.Errorf("test fail %v", err)
	}
	err = masterClient.Connect()
	if err != nil {
		t.Errorf("test fail %v", err)
	}

	time.Sleep(time.Second * 1)

	if !masterClient.isConnect {
		t.Error("test fail")
	}

	// follow client
	servers = []string{followAddress}
	followClient, err := New(servers, timeout, loggerName)
	if err != nil {
		t.Errorf("test fail %v", err)
	}
	err = followClient.Connect()
	if err != nil {
		t.Errorf("test fail %v", err)
	}

	time.Sleep(time.Second * 1)

	if !followClient.isConnect {
		t.Error("test fail")
	}

	// 监控一个不存在的节点
	{
		data := make(chan string, 1)
		path := "/zk_test/123/123"

		// 确定节点确实不存在
		if b, err := masterClient.Exists(path); b || err != nil {
			t.Error("test fail")
		}

		// 开始监控
		go func() {
			err = masterClient.BindWatchValue(path, data)
			if err == nil {
				t.Error("test fail")
			}

			masterClient.UnbindWatchValue(path)
			t.Log("释放监控")
		}()
	}

	// 监控一个存在的节点
	{
		data := make(chan string, 1)
		path := "/zk_test/123"

		// 确认节点存在
		if b, err := masterClient.Exists(path); !b || err != nil {
			t.Error("test fail")
			return
		}

		// 开始监控
		go func() {
			err = masterClient.BindWatchValue(path, data)
			if err != nil {
				t.Errorf("test fail %v", err)
			}
			masterClient.UnbindWatchValue(path)
			t.Log("释放监控")
		}()

		fmt.Println("修改master节点的值")
		actual := <-data
		if !strings.EqualFold(actual, "test") {
			t.Errorf("test fail actual:%s, except:%s", actual, "test")
		}
	}

	// 监控master上的一个节点，然后修改follow对应的节点的值
	{
		data := make(chan string, 1)
		path := "/zk_test/123"

		// 确认节点存在
		if b, err := masterClient.Exists(path); !b || err != nil {
			t.Error("test fail")
			return
		}

		go func() {
			// 确认节点存在
			if b, err := masterClient.Exists(path); !b || err != nil {
				t.Error("test fail")
				return
			}
			// 监控节点值的变化
			err = masterClient.BindWatchValue(path, data)
			if err != nil {
				t.Errorf("test fail %v", err)
			}
			masterClient.UnbindWatchValue(path)
			t.Log("释放监控")
		}()

		fmt.Println("修改follow对应节点的值")
		// 修改follow点对应节点的值
		actual := <-data
		if !strings.EqualFold(actual, "test") {
			t.Errorf("test fail actual:%s, except:%s", actual, "test")
		}
	}

	time.Sleep(time.Second * 1)

	masterClient.Disconnect()
	followClient.Disconnect()
}

// TestUnbindWatchValue 测试取消监控一个节点
func TestUnbindWatchValue(t *testing.T) {
	timeout := time.Second * 1
	loggerName := "zookeeper"
	servers := []string{masterAddress}
	masterClient, err := New(servers, timeout, loggerName)
	if err != nil {
		t.Errorf("test fail %v", err)
	}
	err = masterClient.Connect()
	if err != nil {
		t.Errorf("test fail %v", err)
	}

	time.Sleep(time.Second * 1)

	if !masterClient.isConnect {
		t.Error("test fail")
	}

	// 取消一个没有监控过的节点
	path := "/zk_test/123"
	masterClient.UnbindWatchValue(path)

	// 取消一个路径错误的节点
	path = "home"
	masterClient.UnbindWatchValue(path)
	// 取消一个不存在的节点
	path = "/zk_err_test/err_test"
	masterClient.UnbindWatchValue(path)
	// 取消一个正在监控的节点
	path = "/zk_test/123"
	go func() {
		data := make(chan string, 1)
		err = masterClient.BindWatchValue(path, data)
		if !strings.EqualFold(err.Error(), "zk: zookeeper is closing") {
			t.Errorf("test fail %v", err)
		}
	}()

	masterClient.UnbindWatchValue(path)
	t.Log("释放监控")
	time.Sleep(time.Second * 1)
	masterClient.Disconnect()
<<<<<<< HEAD

=======
	time.Sleep(time.Second * 1)
>>>>>>> 49510abf4584799349d29d0d931d1eacecaa8072
}

// TestBindWatchChildren 测试监控一个节点的子节点
func TestBindWatchChildren(t *testing.T) {
	// master client
	timeout := time.Second * 1
	loggerName := "zookeeper"
	servers := []string{masterAddress}
	masterClient, err := New(servers, timeout, loggerName)
	if err != nil {
		t.Errorf("test fail %v", err)
	}
	err = masterClient.Connect()
	if err != nil {
		t.Errorf("test fail %v", err)
	}

	time.Sleep(time.Second * 1)

	if !masterClient.isConnect {
		t.Error("test fail")
	}

	// 创建一些子节点
	// path := "/zk_test/123/1"
	masterClient.CreateTempNode("/zk_test/123/1", "1")
	masterClient.CreateTempNode("/zk_test/123/2", "2")
	masterClient.CreateTempNode("/zk_test/123/3", "3")

	// 修改子节点的值
	{
		path := "/zk_test/123"
		data := make(chan []string)
		go func() {
			err = masterClient.BindWatchChildren(path, data)
			if err != nil {
				t.Errorf("test fail %v", err)
			}
			masterClient.UnbindWatchChildren(path)
		}()

		// 创建一个节点
		fmt.Println("自动创建一个子节点")
		masterClient.CreateTempNode("/zk_test/123/4", "4")

		actual := <-data
		t.Log(actual)
	}

	// 删除子节点
	{
		path := "/zk_test/123"
		data := make(chan []string)
		go func() {
			err = masterClient.BindWatchChildren(path, data)
			if err != nil {
				t.Errorf("test fail %v", err)
			}
			masterClient.UnbindWatchChildren(path)
		}()

		// 删除子节点
		fmt.Println("手动删除一个子节点")

		actual := <-data
		t.Log(actual)
	}

	// 修改follow对应监控的子节点
	{
		path := "/zk_test/123"
		data := make(chan []string)
		go func() {
			err = masterClient.BindWatchChildren(path, data)
			if err != nil {
				t.Errorf("test fail %v", err)
			}
			masterClient.UnbindWatchChildren(path)
		}()

		fmt.Println("手动修改follow子节点的值")
		actual := <-data
		t.Log(actual)
	}

	masterClient.Disconnect()
	time.Sleep(time.Second * 1)
}

func TestUnbindWatchChildren(t *testing.T) {
	timeout := time.Second * 1
	loggerName := "zookeeper"
	servers := []string{masterAddress}
	masterClient, err := New(servers, timeout, loggerName)
	if err != nil {
		t.Errorf("test fail %v", err)
	}
	err = masterClient.Connect()
	if err != nil {
		t.Errorf("test fail %v", err)
	}

	time.Sleep(time.Second * 1)

	if !masterClient.isConnect {
		t.Error("test fail")
	}

	// 创建一个子节点
	path := "/zk_test/123"
	masterClient.CreateTempNode(path, "")

	// 取消一个没有监控过的节点
	masterClient.UnbindWatchChildren(path)

	// 取消一个监控过的节点
	data := make(chan []string, 1)
	go func() {
		// 创建一些子节点
		masterClient.CreateTempNode(path+"/1", "1")
		masterClient.CreateTempNode(path+"/2", "2")
		masterClient.CreateTempNode(path+"/3", "3")
		err = masterClient.BindWatchChildren(path, data)
		if err == nil {
			t.Errorf("test fail")
		}
	}()
	masterClient.UnbindWatchChildren(path)

	// 关闭连接
	time.Sleep(time.Second * 1)
	masterClient.Disconnect()

	// 启动连接
	masterClient.Connect()

	// 取消一个路径错误的节点
	path = "home"
	masterClient.UnbindWatchChildren(path)

	// 取消一个不存在的节点
	path = "/zk_err_test/err_test"
	masterClient.UnbindWatchChildren(path)

	time.Sleep(time.Second * 1)
	masterClient.Disconnect()
	time.Sleep(time.Second * 1)
}
