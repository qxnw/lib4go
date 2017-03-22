package server

import (
	"golang.org/x/net/context"

	"github.com/qxnw/lib4go/rpc/server/pb"
)

type process struct {
	srv *Server
}

//Request 客户端处理客户端请求
func (r *process) Request(context context.Context, request *pb.RequestContext) (p *pb.ResponseContext, err error) {
	ctx := &Context{}
	ctx.server = r.srv
	ctx.reset("REQUEST", context, request)
	ctx.invoke()
	p = &pb.ResponseContext{}
	p.Status = int32(ctx.Writer.Code)
	p.Result = ctx.Writer.String()
	ctx.Writer.Reset()
	return
}

//Request 客户端处理客户端请求
func (r *process) Query(context context.Context, request *pb.RequestContext) (p *pb.ResponseContext, err error) {
	ctx := &Context{}
	ctx.server = r.srv
	ctx.reset("QUERY", context, request)
	ctx.invoke()
	p = &pb.ResponseContext{}
	p.Status = int32(ctx.Writer.Code)
	p.Result = ctx.Writer.String()
	ctx.Writer.Reset()
	return
}

//Request 客户端处理客户端请求
func (r *process) Update(context context.Context, request *pb.RequestContext) (p *pb.ResponseNoResultContext, err error) {
	ctx := &Context{}
	ctx.server = r.srv
	ctx.reset("UPDATE", context, request)
	ctx.invoke()
	p = &pb.ResponseNoResultContext{}
	p.Status = int32(ctx.Writer.Code)
	ctx.Writer.Reset()
	return
}

//Request 客户端处理客户端请求
func (r *process) Delete(context context.Context, request *pb.RequestContext) (p *pb.ResponseNoResultContext, err error) {
	ctx := &Context{}
	ctx.server = r.srv
	ctx.reset("DELETE", context, request)
	ctx.invoke()
	p = &pb.ResponseNoResultContext{}
	p.Status = int32(ctx.Writer.Code)
	ctx.Writer.Reset()
	return
}

//Request 客户端处理客户端请求
func (r *process) Insert(context context.Context, request *pb.RequestContext) (p *pb.ResponseNoResultContext, err error) {
	ctx := &Context{}
	ctx.server = r.srv
	ctx.reset("INSERT", context, request)
	ctx.invoke()
	p = &pb.ResponseNoResultContext{}
	p.Status = int32(ctx.Writer.Code)
	ctx.Writer.Reset()
	return
}

//Heartbeat 返回心跳数据
func (r *process) Heartbeat(ctx context.Context, in *pb.HBRequest) (*pb.HBResponse, error) {
	return &pb.HBResponse{Pong: in.Ping}, nil
}