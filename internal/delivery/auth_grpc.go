package delivery

import (
	"auth-service/internal/auth"
	"auth-service/internal/grpc/pb"
	"auth-service/internal/grpc/pb/auth_pb"
	"context"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AuthServiceImpl struct {
	auth_pb.UnimplementedAuthServiceServer

	authService auth.AuthService
	logger      *logrus.Logger
}

func NewAuthService(
	authService auth.AuthService,
	logger *logrus.Logger,
) auth_pb.AuthServiceServer {
	return &AuthServiceImpl{
		authService: authService,
		logger:      logger,
	}
}

func (a *AuthServiceImpl) Login(ctx context.Context, request *auth_pb.LoginRequest) (*auth_pb.AuthResponse, error) {
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

func (a *AuthServiceImpl) HasPermission(ctx context.Context, request *auth_pb.CheckPermissionRequest) (*pb.IsSuccess, error) {
	a.logger.Debugf("Does the user have permission %s", request.GetPermission())

	has, err := a.authService.HasPermission(ctx, request.UserID.GetId(), request.GetPermission())
	if err != nil {
		return nil, err
	}

	return &pb.IsSuccess{Value: has}, nil
}

func (a *AuthServiceImpl) RefreshPublicToken(ctx context.Context, request *auth_pb.RefreshTokenRequest) (*auth_pb.AuthResponse, error) {
	//TODO implement me
	panic("implement me")
}
