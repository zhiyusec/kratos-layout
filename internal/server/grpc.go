package server

import (
	v1 "github.com/zhiyusec/kratos-layout/api/todo/v1"
	"github.com/zhiyusec/kratos-layout/internal/conf"
	"github.com/zhiyusec/kratos-layout/internal/service"

	"github.com/go-kratos/kratos/v3/middleware/recovery"
	"github.com/go-kratos/kratos/v3/transport/grpc"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Server, todo *service.TodoService) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
		),
	}
	if c.GRPC.Network != "" {
		opts = append(opts, grpc.Network(c.GRPC.Network))
	}
	if c.GRPC.Addr != "" {
		opts = append(opts, grpc.Address(c.GRPC.Addr))
	}
	if c.GRPC.Timeout.AsDuration() != 0 {
		opts = append(opts, grpc.Timeout(c.GRPC.Timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)
	v1.RegisterTodoServiceServer(srv, todo)
	return srv
}
