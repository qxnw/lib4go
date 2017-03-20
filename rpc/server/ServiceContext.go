package server

import (
	"golang.org/x/net/context"

	"github.com/qxnw/lib4go/rpc/server/pb"
)

type serverCaller struct {
	srv *Server
}

//Request 客户端处理客户端请求
func (r *serverCaller) Request(context context.Context, request *pb.RequestContext) (p *pb.ResponseContext, err error) {
	s, d, err := r.srv.callback(request.Session, request.Sevice, request.Input)
	if err != nil {
		return
	}
	p = &pb.ResponseContext{Status: int32(s), Result: d}
	return
}

//Heartbeat 返回心跳数据
func (r *serverCaller) Heartbeat(ctx context.Context, in *pb.HBRequest) (*pb.HBResponse, error) {
	return &pb.HBResponse{Pong: in.Ping}, nil
}
