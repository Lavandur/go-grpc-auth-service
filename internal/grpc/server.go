package grpc

import (
	"auth-service/internal/grpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func NewGRPCServer(authService pb.AuthServiceServer) *grpc.Server {
	grpcServer := grpc.NewServer()

	reflection.Register(grpcServer)
	pb.RegisterAuthServiceServer(grpcServer, authService)

	return grpcServer
}
