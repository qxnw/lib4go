package server

import (
	"net"

	"github.com/qxnw/lib4go/rpc/server/pb"

	"google.golang.org/grpc"
)

//Server RPC Server
type Server struct {
	server   *grpc.Server
	address  string
	callback func(string, string, string) (int, string, error)
}

//NewServer 初始化
func NewServer(f func(string, string, string) (int, string, error)) *Server {
	return &Server{callback: f}
}

//Start 启动RPC服务器
func (r *Server) Start(address string) (err error) {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return
	}
	r.server = grpc.NewServer()
	pb.RegisterARSServer(r.server, &process{srv: r})
	r.server.Serve(lis)
	return
}

//Close 关闭连接
func (r *Server) Close() {
	if r.server != nil {
		r.server.GracefulStop()
	}
}
