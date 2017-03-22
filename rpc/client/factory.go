package client

import (
	"time"

	"github.com/qxnw/lib4go/concurrent/cmap"
	"github.com/qxnw/lib4go/rpc/client/balancer"
	"google.golang.org/grpc/naming"
)

//RPCClientFactory rpc client factory
type RPCClientFactory struct {
	cache   cmap.ConcurrentMap
	address string
	opts    []ClientOption

	lb balancer.CustomerBalancer
	*facotryOption
}

type facotryOption struct {
	logger       Logger
	timerout     time.Duration
	balancerType int
	servers      string
	local        string
	mode         int
}

const (
	ZKBalancer = iota + 1
	FileBalancer
)
const (
	RoundRobin = iota
	LocalFirst
)

//FactoryOption 客户端配置选项
type FactoryOption func(*facotryOption)

//WithFactoryLogger 设置日志记录器
func WithFactoryLogger(log Logger) FactoryOption {
	return func(o *facotryOption) {
		o.logger = log
	}
}

//WithRoundRobin 设置为轮询负载
func WithRoundRobin() FactoryOption {
	return func(o *facotryOption) {
		o.mode = RoundRobin
	}
}

//WithLocalFirst 设置为轮询负载
func WithLocalFirst(local string) FactoryOption {
	return func(o *facotryOption) {
		o.mode = LocalFirst
		o.local = local
	}
}

//WithZKBalancer 设置基于Zookeeper服务发现的负载均衡器
func WithZKBalancer(servers string, timeout time.Duration) FactoryOption {
	return func(o *facotryOption) {
		o.servers = servers
		o.timerout = timeout
		o.balancerType = ZKBalancer
	}
}

//WithFileBalancer 设置本地优先负载均衡器
func WithFileBalancer(f string) FactoryOption {
	return func(o *facotryOption) {
		o.servers = f
		o.balancerType = FileBalancer
	}
}

//NewRPCClientFactory new rpc client factory
func NewRPCClientFactory(address string, opts ...FactoryOption) (f *RPCClientFactory) {
	f = &RPCClientFactory{
		address:       address,
		cache:         cmap.New(),
		facotryOption: &facotryOption{},
	}
	for _, opt := range opts {
		opt(f.facotryOption)
	}
	return
}

//Get 获取rpc client
func (r *RPCClientFactory) Get(service string) (c *Client, err error) {
	_, client, err := r.cache.SetIfAbsentCb(service, func(i ...interface{}) (interface{}, error) {
		opts := make([]ClientOption, 0, 0)
		opts = append(opts, WithLogger(r.logger))
		if r.balancerType > 0 {
			var rs naming.Resolver
			switch r.balancerType {
			case ZKBalancer:
				rs = balancer.NewZKResolver(service, r.local, time.Second)
			}
			if rs != nil {
				switch r.mode {
				case RoundRobin:
					opts = append(opts, WithRoundRobinBalancer(rs, service, time.Second, map[string]int{}))
				case LocalFirst:
					opts = append(opts, WithLocalFirstBalancer(rs, service, r.local, map[string]int{}))
				}
			}
		}
		return NewClient(r.address, opts...), nil
	})
	if err != nil {
		return
	}
	c = client.(*Client)
	return
}
