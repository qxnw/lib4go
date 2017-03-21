package client

import (
	"fmt"
	"io"
	"time"

	"github.com/lunny/log"
	"github.com/qxnw/lib4go/rpc/client/balancer"
	"github.com/qxnw/lib4go/rpc/client/balancer/zkbalancer"
	"github.com/qxnw/lib4go/rpc/server/pb"

	"os"

	"strings"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
)

//Logger 日志组件
type Logger interface {
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Fatalln(args ...interface{})
	Print(args ...interface{})
	Printf(format string, args ...interface{})
	Println(args ...interface{})
}

//Client client
type Client struct {
	address       string
	conn          *grpc.ClientConn
	opts          *clientOption
	client        pb.ARSClient
	longTicker    *time.Ticker
	lastRequest   time.Time
	hasRunChecker bool
	IsConnect     bool
	isClose       bool
}

type clientOption struct {
	connectionTimeout time.Duration
	log               Logger
	balancer          grpc.Balancer
	serviceGroup      string
}

//ClientOption 客户端配置选项
type ClientOption func(*clientOption)

//WithConnectionTimeout 网络连接超时时长
func WithConnectionTimeout(t time.Duration) ClientOption {
	return func(o *clientOption) {
		o.connectionTimeout = t
	}
}

//WithLogger 设置日志记录器
func WithLogger(log Logger) ClientOption {
	return func(o *clientOption) {
		o.log = log
	}
}

//WithZKRoundRobinBalancer 设置负载均衡器
func WithZKRoundRobinBalancer(serviceGroup string, timeout time.Duration) ClientOption {
	return func(o *clientOption) {
		r := zkbalancer.NewResolver(serviceGroup, timeout)
		o.serviceGroup = serviceGroup
		o.balancer = balancer.RoundRobin(r)
	}
}

//NewClient 创建客户端
func NewClient(address string, opts ...ClientOption) *Client {
	client := &Client{address: address, opts: &clientOption{connectionTimeout: time.Second * 3}}
	for _, opt := range opts {
		opt(client.opts)
	}
	if client.opts.log == nil {
		client.opts.log = NewLogger(os.Stdout)
	}
	grpclog.SetLogger(client.opts.log)
	client.connect()
	return client
}

//Connect 连接服务器，如果当前无法连接系统会定时自动重连
func (c *Client) connect() (b bool) {
	if c.IsConnect {
		return
	}
	var err error
	if c.opts.balancer == nil {
		c.conn, err = grpc.Dial(c.address, grpc.WithInsecure(), grpc.WithTimeout(c.opts.connectionTimeout))
	} else {
		ctx, _ := context.WithTimeout(context.Background(), c.opts.connectionTimeout)
		c.conn, err = grpc.DialContext(ctx, c.address, grpc.WithInsecure(), grpc.WithBalancer(c.opts.balancer))
	}
	if err != nil {
		c.IsConnect = false
		return c.IsConnect
	}
	c.client = pb.NewARSClient(c.conn)
	//检查是否已连接到服务器
	response, er := c.client.Heartbeat(context.Background(), &pb.HBRequest{Ping: 0})
	c.IsConnect = er == nil && response.Pong == 0
	return c.IsConnect
}

//Request 发送请求
func (c *Client) Request(session string, service string, data string, kv ...string) (status int, result string, err error) {
	c.lastRequest = time.Now()
	if !strings.HasPrefix(service, c.opts.serviceGroup) {
		return 500, "", fmt.Errorf("服务:%s调用失败", service)
	}
	response, err := c.client.Request(metadata.NewContext(context.Background(), metadata.Pairs(kv...)), &pb.RequestContext{Session: session, Sevice: service, Input: data},
		grpc.FailFast(true))
	if err != nil {
		c.IsConnect = false
		return
	}
	status = int(response.Status)
	result = response.GetResult()
	c.IsConnect = true
	return
}

//logInfof 日志记录
func (c *Client) logInfof(format string, msg ...interface{}) {
	if c.opts.log == nil {
		return
	}
	c.opts.log.Printf(format, msg...)
}

//Close 关闭连接
func (c *Client) Close() {
	c.isClose = true
	if c.longTicker != nil {
		c.longTicker.Stop()
	}
	if c.conn != nil {
		c.conn.Close()
	}
}

//NewLogger 创建日志组件
func NewLogger(out io.Writer) Logger {
	l := log.New(out, "[grpc.client] ", log.Ldefault())
	l.SetOutputLevel(log.Ldebug)
	return &nLogger{Logger: l}
}

type nLogger struct {
	*log.Logger
}

func (n *nLogger) Fatalln(args ...interface{}) {
	n.Fatal(args...)
}
