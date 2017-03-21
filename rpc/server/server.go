package server

import (
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"os"

	"github.com/qxnw/lib4go/rpc/server/pb"
	"google.golang.org/grpc"
)

//Server RPC Server
type Server struct {
	server     *grpc.Server
	serverName string
	address    string
	process    *process
	ctxPool    sync.Pool
	ErrHandler Handler
	*serverOption
	Router
}

//Version 获取当前版本号
func Version() string {
	return "0.0.1"
}

type serverOption struct {
	logger   Logger
	handlers []Handler
}

//Option 配置选项
type Option func(*serverOption)

//WithLogger 设置日志记录组件
func WithLogger(logger Logger) Option {
	return func(o *serverOption) {
		o.logger = logger
	}
}

//WithInfluxMetric 设置基于influxdb的系统监控组件
func WithInfluxMetric(host string, dataBase string, userName string, password string, timeSpan time.Duration) Option {
	return func(o *serverOption) {
		o.handlers = append(o.handlers, NewInfluxMetric(host, dataBase, userName, password, timeSpan))
	}
}

//WithHandlers 添加插件
func WithHandlers(handlers ...Handler) Option {
	return func(o *serverOption) {
		o.handlers = append(o.handlers, handlers...)
	}
}

var (
	//ClassicHandlers 标准插件
	ClassicHandlers = []Handler{
		Logging(),
		Recovery(false),
		Return(),
		Param(),
		Contexts(),
	}
)

//NewServer 初始化
func NewServer(address string, opts ...Option) *Server {
	s := &Server{address: GetAddress(address), Router: NewRouter()}
	s.serverOption = &serverOption{}
	s.logger = NewLogger(os.Stdout)
	s.process = &process{srv: s}
	s.ErrHandler = Errors()
	for _, opt := range opts {
		opt(s.serverOption)
	}
	s.Use(ClassicHandlers...)
	return s
}

//Use 使用新的插件
func (s *Server) Use(handlers ...Handler) {
	s.handlers = append(s.handlers, handlers...)
}

//Run 启动RPC服务器
func (s *Server) Run() (err error) {
	s.logger.Info("Listening on " + s.address)
	lis, err := net.Listen("tcp", s.address)
	if err != nil {
		return
	}
	s.server = grpc.NewServer()
	pb.RegisterRPCServer(s.server, s.process)
	s.server.Serve(lis)
	return
}

//Close 关闭连接
func (s *Server) Close() {
	if s.server != nil {
		s.server.GracefulStop()
	}
}

//Logger 获取日志组件
func (s *Server) Logger() Logger {
	return s.logger
}

//Request 设置Request路由
func (s *Server) Request(service string, c interface{}, middlewares ...Handler) {
	s.Route([]string{"REQUEST"}, service, c, middlewares...)
}

//Query 设置Query路由
func (s *Server) Query(service string, c interface{}, middlewares ...Handler) {
	s.Route([]string{"QUERY"}, service, c, middlewares...)
}

//Insert 设置Insert路由
func (s *Server) Insert(service string, c interface{}, middlewares ...Handler) {
	s.Route([]string{"INSERT"}, service, c, middlewares...)
}

//Delete 设置Delete路由
func (s *Server) Delete(service string, c interface{}, middlewares ...Handler) {
	s.Route([]string{"DELETE"}, service, c, middlewares...)
}

//Update 设置Update路由
func (s *Server) Update(service string, c interface{}, middlewares ...Handler) {
	s.Route([]string{"UPDATE"}, service, c, middlewares...)
}
func GetAddress(args ...interface{}) string {
	var host string
	var port int

	if len(args) == 1 {
		switch arg := args[0].(type) {
		case string:
			addrs := strings.Split(args[0].(string), ":")
			if len(addrs) == 1 {
				host = addrs[0]
			} else if len(addrs) >= 2 {
				host = addrs[0]
				_port, _ := strconv.ParseInt(addrs[1], 10, 0)
				port = int(_port)
			}
		case int:
			port = arg
		}
	} else if len(args) >= 2 {
		if arg, ok := args[0].(string); ok {
			host = arg
		}
		if arg, ok := args[1].(int); ok {
			port = arg
		}
	}

	if len(host) == 0 {
		host = "0.0.0.0"
	}
	if port == 0 {
		port = 8000
	}

	addr := host + ":" + strconv.FormatInt(int64(port), 10)

	return addr
}
