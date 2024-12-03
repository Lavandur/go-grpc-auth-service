package delivery

import (
	"auth-service/internal/grpc/pb/auth_pb"
	"auth-service/internal/grpc/pb/roles_pb"
	"auth-service/internal/grpc/pb/users_pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func NewGRPCServer(
	authService auth_pb.AuthServiceServer,
	userService users_pb.UserServiceServer,
	roleService roles_pb.RoleServiceServer,
) *grpc.Server {
	grpcServer := grpc.NewServer()

	reflection.Register(grpcServer)
	auth_pb.RegisterAuthServiceServer(grpcServer, authService)
	users_pb.RegisterUserServiceServer(grpcServer, userService)
	roles_pb.RegisterRoleServiceServer(grpcServer, roleService)

	return grpcServer
}
