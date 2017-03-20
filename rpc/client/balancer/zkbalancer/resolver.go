package zkbalancer

import (
	"fmt"
	"strings"
	"time"

	"github.com/qxnw/lib4go/zk"

	"google.golang.org/grpc/naming"
)

type resolver struct {
	timeout     time.Duration
	serviceRoot string
	prefix      string
}

//NewResolver 返回服务解析器
func NewResolver(serviceRoot string, prefix string, timeout time.Duration) *resolver {
	return &resolver{timeout: timeout, serviceRoot: serviceRoot, prefix: prefix}
}

// Resolve to resolve the service from zookeeper, target is the dial address of zookeeper
// target example: "192.168.0.159:2181;192.168.0.154:2181"
func (re *resolver) Resolve(target string) (naming.Watcher, error) {
	client, err := zk.New(strings.Split(target, ";"), re.timeout)
	if err != nil {
		return nil, fmt.Errorf("grpclb: creat zookeeper client failed: %s", err.Error())
	}
	return &watcher{re: re, client: client}, nil
}
