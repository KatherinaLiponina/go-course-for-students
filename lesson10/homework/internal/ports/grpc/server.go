package grpc

import (
	context "context"
	"homework10/internal/app"
	"log"
	"time"

	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

func NewGRPCServer(port string, a app.App) *grpc.Server {
	customFunc := func(p interface{}) (err error) {
		return status.Errorf(codes.Unknown, "panic triggered: %v", p)
	}
	opts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(customFunc),
	}
	service := &AdUserService{App: a}
	server := grpc.NewServer(grpc.ChainUnaryInterceptor(UnaryServerInterceptor),
		grpc.ChainUnaryInterceptor(grpc_recovery.UnaryServerInterceptor(opts...)))
	RegisterAdServiceServer(server, service)
	return server
}

func UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Println(time.Now().GoString() + ": " + info.FullMethod)

	return handler(ctx, req)
}
