package zk

import (
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/qxnw/lib4go/encoding"
)

var address = "192.168.0.166:2181"

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
	path := "/err_path/err_node"
	err = client.Delete(path)
	if err == nil {
		t.Error("test fail")
	}

	// 删除一个存在的节点
	path = "/zk_test/123"
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

	// 删除一个节点存在，有子节点
	path = "zk_test"
	err = client.Delete(path)
	if err == nil {
		t.Error("test fail")
	}

	client.Disconnect()
}

// TestGetPaths 获取当前路径的子路径
func TestGetPaths(t *testing.T) {
	client := &ZookeeperClient{}

	// 获取当前路径的子路径
	path := "/home/champly/test"
	paths := client.getPaths(path)
	if len(paths) != 3 {
		t.Errorf("test fail actual : %d", len(paths))
	}
	if !strings.EqualFold(paths[0], "/home") || !strings.EqualFold(paths[1], "/home/champly") || !strings.EqualFold(paths[2], "/home/champly/test") {
		t.Log(paths)
		t.Error("test fail")
	}

	// 结尾是/
	path = "/home/champly/test/"
	paths = client.getPaths(path)
	if len(paths) != 4 {
		t.Errorf("test fail actual : %d", len(paths))
	}
	if !strings.EqualFold(paths[0], "/home") || !strings.EqualFold(paths[1], "/home/champly") || !strings.EqualFold(paths[2], "/home/champly/test") || !strings.EqualFold(paths[3], "/home/champly/test/") {
		t.Log(paths)
		t.Error("test fail")
	}

	// 不是以 / 开头
	path = "home/champly/test"
	paths = client.getPaths(path)
	if len(paths) != 2 {
		t.Errorf("test fail actual : %d", len(paths))
	}
	if !strings.EqualFold(paths[0], "/champly") || !strings.EqualFold(paths[1], "/champly/test") {
		t.Log(paths)
		t.Error("test fail")
	}

	// 路径包含特殊字符
	path = "\\\\123/#!@_"
	paths = client.getPaths(path)
	if len(paths) != 1 {
		t.Errorf("test fail actual : %d", len(paths))
	}
	if !strings.EqualFold(paths[0], "/#!@_") {
		t.Log(paths)
		t.Error("test fail")
	}
}

// TestGetDir 获取当前路径的目录
func TestGetDir(t *testing.T) {
	client := &ZookeeperClient{}

	// 当前路径为子节点路径
	path := "/home/champly/test"
	dir := client.GetDir(path)
	if !strings.EqualFold(dir, "/home/champly") {
		t.Errorf("test fail actual:%s, except:%s", dir, "/home/champly")
	}

	// 当前路径为父节点
	path = "/home/champly/test/"
	dir = client.GetDir(path)
	if !strings.EqualFold(dir, "/home/champly/test") {
		t.Errorf("test fail actual:%s, except:%s", dir, "/home/champly/test")
	}
}

// TestCreatePersistentNode 测试创建持久化的节点
func TestCreatePersistentNode(t *testing.T) {
	timeout := time.Second * 1
	loggerName := "zookeeper"
	servers := []string{address}
	client, err := New(servers, timeout, loggerName)
	if err != nil {
		t.Errorf("test fail %v", err)
	}

	// 创建一个持久化的节点
	{
		err = client.Connect()
		if err != nil {
			t.Errorf("test fail %v", err)
		}

		time.Sleep(time.Second * 1)

		if !client.isConnect {
			t.Error("test fail")
		}

		path := "/zk_test/123"
		err = client.CreatePersistentNode(path, "")
		if err != nil {
			t.Errorf("test fail : %v", err)
		}

		// 校验断开连接之后节点是否存在
		client.Disconnect()
		err = client.Connect()
		if err != nil {
			t.Errorf("test fail %v", err)
		}

		time.Sleep(time.Second * 1)

		if !client.isConnect {
			t.Error("test fail")
		}
		b, err := client.Exists(path)
		if err != nil {
			t.Errorf("create persistent node fail : %v", err)
		}
		if !b {
			t.Error("create persistent node fail")
		}

		client.Disconnect()
	}

	// 创建多级目录的节点
	{
		err = client.Connect()
		if err != nil {
			t.Errorf("test fail %v", err)
		}

		time.Sleep(time.Second * 1)

		if !client.isConnect {
			t.Error("test fail")
		}

		path := "/zk_test1/123/123"
		err = client.CreatePersistentNode(path, "")
		if err != nil {
			t.Errorf("test fail : %v", err)
		}

		// 校验断开连接之后节点是否存在
		client.Disconnect()
		err = client.Connect()
		if err != nil {
			t.Errorf("test fail %v", err)
		}

		time.Sleep(time.Second * 1)

		if !client.isConnect {
			t.Error("test fail")
		}
		b, err := client.Exists(path)
		if err != nil {
			t.Errorf("create persistent node fail : %v", err)
		}
		if !b {
			t.Error("create persistent node fail")
		}

		// // 删除节点，方便下次测试
		// paths := client.getPaths(path)
		// for i := len(paths) - 1; i >= 0; i-- {
		// 	err = client.Delete(paths[i])
		// 	if err != nil {
		// 		t.Errorf("delete node fail : %s", paths[i])
		// 	}
		// }

		client.Disconnect()
	}

	// 创建重复的节点
	{
		err = client.Connect()
		if err != nil {
			t.Errorf("test fail %v", err)
		}

		time.Sleep(time.Second * 1)

		if !client.isConnect {
			t.Error("test fail")
		}

		path := "/zk_test1/123/123"
		err = client.CreatePersistentNode(path, "")
		if err != nil {
			t.Errorf("test fail : %v", err)
		}

		// 校验断开连接之后节点是否存在
		client.Disconnect()
		err = client.Connect()
		if err != nil {
			t.Errorf("test fail %v", err)
		}

		time.Sleep(time.Second * 1)

		if !client.isConnect {
			t.Error("test fail")
		}
		b, err := client.Exists(path)
		if err != nil {
			t.Errorf("create persistent node fail : %v", err)
		}
		if !b {
			t.Error("create persistent node fail")
		}

		// 删除节点，方便下次测试
		paths := client.getPaths(path)
		for i := len(paths) - 1; i >= 0; i-- {
			err = client.Delete(paths[i])
			if err != nil {
				t.Errorf("delete node fail : %s", paths[i])
			}
		}

		client.Disconnect()
	}
}

// TestTempNode 测试创建临时节点
func TestTempNode(t *testing.T) {
	timeout := time.Second * 1
	loggerName := "zookeeper"
	servers := []string{address}
	client, err := New(servers, timeout, loggerName)
	if err != nil {
		t.Errorf("test fail %v", err)
	}

	// 创建一个临时节点，父节点不存在
	{
		err = client.Connect()
		if err != nil {
			t.Errorf("test fail %v", err)
		}

		time.Sleep(time.Second * 1)

		if !client.isConnect {
			t.Error("test fail")
		}

		path := "/zk_test1/123/123"
		rpath, err := client.CreateTempNode(path, "")
		if err != nil {
			t.Errorf("test fail : %v", err)
		}

		// 判断节点是否存在
		b, err := client.Exists(rpath)
		if err != nil {
			t.Errorf("create persistent node fail : %v", err)
		}
		if !b {
			t.Error("create persistent node fail")
		}

		// 校验断开连接之后节点是否存在
		client.Disconnect()
		err = client.Connect()
		if err != nil {
			t.Errorf("test fail %v", err)
		}

		time.Sleep(time.Second * 1)

		if !client.isConnect {
			t.Error("test fail")
		}
		b, err = client.Exists(rpath)
		if err != nil {
			t.Errorf("create temp node fail : %v", err)
		}
		if b {
			t.Error("create temp node fail")
		}

		// 判断父节点是否存在
		dir := client.GetDir(rpath)
		b, err = client.Exists(dir)
		if err != nil {
			t.Errorf("create temp node fail : %v", err)
		}
		if !b {
			t.Error("create temp node fail")
		}

		client.Disconnect()
	}

	// 创建一个存在的临时节点
	{
		err = client.Connect()
		if err != nil {
			t.Errorf("test fail %v", err)
		}

		time.Sleep(time.Second * 1)

		if !client.isConnect {
			t.Error("test fail")
		}

		// 创建第一个临时节点
		path := "/zk_test1/123/123"
		rpath, err := client.CreateTempNode(path, "")
		if err != nil {
			t.Errorf("test fail : %v", err)
		}

		// 创建第二个临时节点
		path = "/zk_test1/123/123"
		_, err = client.CreateTempNode(path, "")
		if err == nil {
			t.Error("test fail")
		}

		// 判断节点是否存在
		b, err := client.Exists(rpath)
		if err != nil {
			t.Errorf("create persistent node fail : %v", err)
		}
		if !b {
			t.Error("create persistent node fail")
		}

		// 校验断开连接之后节点是否存在
		client.Disconnect()
		err = client.Connect()
		if err != nil {
			t.Errorf("test fail %v", err)
		}

		time.Sleep(time.Second * 1)

		if !client.isConnect {
			t.Error("test fail")
		}
		b, err = client.Exists(rpath)
		if err != nil {
			t.Errorf("create temp node fail : %v", err)
		}
		if b {
			t.Error("create temp node fail")
		}

		// 判断父节点是否存在
		dir := client.GetDir(rpath)
		b, err = client.Exists(dir)
		if err != nil {
			t.Errorf("create temp node fail : %v", err)
		}
		if !b {
			t.Error("create temp node fail")
		}

		client.Disconnect()
	}

	// 创建一个已存在的持久化节点
	{
		err = client.Connect()
		if err != nil {
			t.Errorf("test fail %v", err)
		}

		time.Sleep(time.Second * 1)

		if !client.isConnect {
			t.Error("test fail")
		}

		path := "/zk_test1/123"
		_, err := client.CreateTempNode(path, "")
		if err == nil {
			t.Errorf("test fail : %v", err)
		}

		// 删除节点，方便下次测试
		paths := client.getPaths(path)
		for i := len(paths) - 1; i >= 0; i-- {
			err = client.Delete(paths[i])
			if err != nil {
				t.Errorf("delete node fail : %s", paths[i])
			}
		}

		client.Disconnect()
	}

	// 节点名为 /
	{
		err = client.Connect()
		if err != nil {
			t.Errorf("test fail %v", err)
		}

		time.Sleep(time.Second * 1)

		if !client.isConnect {
			t.Error("test fail")
		}

		path := "/"
		_, err := client.CreateTempNode(path, "")
		if err == nil {
			t.Errorf("test fail : %v", err)
		}

		client.Disconnect()
	}
}

// TestCreateSeqNode 测试创建有临时序节点
func TestCreateSeqNode(t *testing.T) {
	timeout := time.Second * 1
	loggerName := "zookeeper"
	servers := []string{address}
	client, err := New(servers, timeout, loggerName)
	if err != nil {
		t.Errorf("test fail %v", err)
	}

	// 创建一个临时有序节点
	{
		err = client.Connect()
		if err != nil {
			t.Errorf("test fail %v", err)
		}

		time.Sleep(time.Second * 1)

		if !client.isConnect {
			t.Error("test fail")
		}

		path := "/zk_test1/123/test"
		rpath1, err := client.CreateSeqNode(path, "")
		if err != nil {
			t.Errorf("test fail : %v", err)
		}

		rpath2, err := client.CreateSeqNode(path, "")
		if err != nil {
			t.Errorf("test fail : %v", err)
		}

		t.Log(rpath1)
		t.Log(rpath2)
		// 校验节点的有序性
		r1, err := strconv.Atoi(string(rpath1[len(rpath1)-4 : len(rpath1)]))
		if err != nil {
			t.Errorf("test fail : %v", err)
		}
		r2, err := strconv.Atoi(string(rpath2[len(rpath2)-4 : len(rpath2)]))
		if err != nil {
			t.Errorf("test fail : %v", err)
		}

		t.Logf("r1 : %d, r2 : %d", r1, r2)
		if r2-r1 != 1 {
			t.Errorf("test fail")
		}

		// 判断节点是否存在
		b, err := client.Exists(rpath1)
		if err != nil {
			t.Errorf("create persistent node fail : %v", err)
		}
		if !b {
			t.Error("create persistent node fail")
		}
		b, err = client.Exists(rpath2)
		if err != nil {
			t.Errorf("create persistent node fail : %v", err)
		}
		if !b {
			t.Error("create persistent node fail")
		}

		// 校验断开连接之后节点是否存在
		client.Disconnect()
		err = client.Connect()
		if err != nil {
			t.Errorf("test fail %v", err)
		}

		time.Sleep(time.Second * 1)

		if !client.isConnect {
			t.Error("test fail")
		}
		b, err = client.Exists(rpath1)
		if err != nil {
			t.Errorf("create temp node fail : %v", err)
		}
		if b {
			t.Error("create temp node fail")
		}
		b, err = client.Exists(rpath2)
		if err != nil {
			t.Errorf("create temp node fail : %v", err)
		}
		if b {
			t.Error("create temp node fail")
		}

		// 判断父节点是否存在
		dir := client.GetDir(rpath1)
		b, err = client.Exists(dir)
		if err != nil {
			t.Errorf("create temp node fail : %v", err)
		}
		if !b {
			t.Error("create temp node fail")
		}

		client.Disconnect()
	}
}

func TestGetValue(t *testing.T) {
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

	// 获取永久节点的值
	{
		inputs := map[string]string{
			"/zk_test1/123/1": "Tom1",
			"/zk_test1/123/2": "!@#!!#……",
			"/zk_test1/123/3": "",
			"/zk_test1/123/4": " ",
			"/zk_test1/123/5": "你好",
		}

		for path, except := range inputs {
			err = client.CreatePersistentNode(path, except)
			if err != nil {
				t.Errorf("test fail : %v", err)
			}

			except, err = encoding.Convert([]byte(except), "gbk")
			if err != nil {
				client.Delete(path)
				continue
			}

			// 获取值
			actual, err := client.GetValue(path)
			if err != nil {
				t.Errorf("test fail %v", err)
			}

			if !strings.EqualFold(actual, except) {
				t.Errorf("test fail , actual:%s, except:%s", actual, except)
			}

			// 删除节点
			client.Delete(path)
		}
	}

	// 获取临时节点的值
	{
		inputs := map[string]string{
			"/zk_test1/123/1": "Tom1",
			"/zk_test1/123/2": "!@#!!#……",
			"/zk_test1/123/3": "",
			"/zk_test1/123/4": " ",
			"/zk_test1/123/5": "你好",
		}

		for path, except := range inputs {
			path, err = client.CreateTempNode(path, except)
			if err != nil {
				t.Errorf("test fail : %v", err)
			}

			except, err = encoding.Convert([]byte(except), "gbk")
			if err != nil {
				client.Delete(path)
				continue
			}

			// 获取值
			actual, err := client.GetValue(path)
			if err != nil {
				t.Errorf("test fail %v", err)
			}

			if !strings.EqualFold(actual, except) {
				t.Errorf("test fail , actual:%s, except:%s", actual, except)
			}

			// 删除节点
			client.Delete(path)
		}
	}

	// 获取临时有序节点的值
	{
		inputs := map[string]string{
			"/zk_test1/123/1": "Tom1",
			"/zk_test1/123/2": "!@#!!#……",
			"/zk_test1/123/3": "",
			"/zk_test1/123/4": " ",
			"/zk_test1/123/5": "你好",
		}

		for path, except := range inputs {
			path, err = client.CreateSeqNode(path, except)
			if err != nil {
				t.Errorf("test fail : %v", err)
			}

			except, err = encoding.Convert([]byte(except), "gbk")
			if err != nil {
				client.Delete(path)
				continue
			}

			// 获取值
			actual, err := client.GetValue(path)
			if err != nil {
				t.Errorf("test fail %v", err)
			}

			if !strings.EqualFold(actual, except) {
				t.Errorf("test fail , actual:%s, except:%s", actual, except)
			}

			// 删除节点
			client.Delete(path)
		}
	}

	// 获取一个不存在节点的值
	{
		_, err = client.GetValue("/zk_test_err/no_node")
		if err == nil {
			t.Error("test fail")
		}
	}

	path := "/zk_test1/123"
	for _, p := range client.getPaths(path) {
		client.Delete(p)
	}
	client.Disconnect()
}

// TestGetChildren 测试获取当前路径下的所有子节点
func TestGetChildren(t *testing.T) {
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

	// 节点存在，有子节点
	path := "/zk_test"
	paths, err := client.GetChildren(path)
	if err != nil {
		t.Errorf("test fail %v", err)
	}
	t.Log(paths)
	if len(paths) == 0 {
		t.Error("test fail")
	}

	// 节点不存在
	path = "/zk_test/test"
	paths, err = client.GetChildren(path)
	t.Log(paths)
	if err != nil {
		t.Errorf("test fail %v", err)
	}
	if len(paths) != 0 {
		t.Errorf("test fail")
	}

	// 没有子节点
	path = "/zk_test/123"
	paths, err = client.GetChildren(path)
	t.Log(paths)
	if err != nil {
		t.Errorf("test fail %v", err)
	}
	if len(paths) != 0 {
		t.Errorf("test fail")
	}

	client.Disconnect()
}
