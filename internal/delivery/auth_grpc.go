package delivery

import (
	"auth-service/internal/auth"
	"auth-service/internal/grpc/pb"
	"auth-service/internal/grpc/pb/auth_pb"
	"context"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AuthGRPCService struct {
	auth_pb.UnimplementedAuthServiceServer

	authService auth.AuthService
	logger      *logrus.Logger
}

func NewAuthService(
	authService auth.AuthService,
	logger *logrus.Logger,
) auth_pb.AuthServiceServer {
	return &AuthGRPCService{
		authService: authService,
		logger:      logger,
	}
}

func (a *AuthGRPCService) Login(ctx context.Context, request *auth_pb.LoginRequest) (*auth_pb.AuthResponse, error) {
	a.logger.Debugf("Login user with login %s", request.GetLogin())

	response, err := a.authService.Login(ctx, request.GetLogin(), request.GetPassword())
	if err != nil {
		return nil, err
	}

	return &auth_pb.AuthResponse{
		AccessToken: response.PublicToken,
		ExpiresAt:   timestamppb.New(response.PublicTokenExpiry),
	}, nil
}

func (a *AuthGRPCService) HasPermission(ctx context.Context, request *auth_pb.CheckPermissionRequest) (*pb.IsSuccess, error) {
	a.logger.Debugf("Does the user have permission %s", request.GetPermission())

	has := a.authService.HasPermission(ctx, request.GetPermission())

	return &pb.IsSuccess{Value: has}, nil
}
