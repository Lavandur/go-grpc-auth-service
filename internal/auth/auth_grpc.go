package auth

import (
	"auth-service/internal/grpc/pb/auth_pb"
	"context"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AuthServiceImpl struct {
	auth_pb.UnimplementedAuthServiceServer

	authService AuthService
	logger      *logrus.Logger
}

func NewAuthService(
	authService AuthService,
	logger *logrus.Logger,
) auth_pb.AuthServiceServer {
	return &AuthServiceImpl{
		authService: authService,
		logger:      logger,
	}
}

func (a *AuthServiceImpl) Login(ctx context.Context, request *auth_pb.LoginRequest) (*auth_pb.AuthResponse, error) {
	a.logger.Debugf("Login user with login %s", request.Login)

	response, err := a.authService.Login(ctx, request.GetLogin(), request.GetPassword())
	if err != nil {
		return nil, err
	}

	return &auth_pb.AuthResponse{
		AccessToken: response.PublicToken,
		ExpiresAt:   timestamppb.New(response.PublicTokenExpiry),
	}, nil
}

func (a *AuthServiceImpl) RefreshPublicToken(ctx context.Context, request *auth_pb.RefreshTokenRequest) (*auth_pb.AuthResponse, error) {
	//TODO implement me
	panic("implement me")
}
