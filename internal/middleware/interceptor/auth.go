package interceptor

import (
	"auth-service/internal/auth/auth_service"
	"context"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"slices"
)

type AuthInterceptor struct {
	paseto auth_service.PasetoAuth
	logger *logrus.Logger
}

func NewAuthInterceptor(
	paseto auth_service.PasetoAuth,
	logger *logrus.Logger,
) AuthInterceptor {
	return AuthInterceptor{
		paseto: paseto,
		logger: logger,
	}
}

func (a *AuthInterceptor) AuthUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		a.logger.Debugf("Executing request with name: %s", info.FullMethod)

		withoutAuth := []string{"/auth_pb.AuthService/login", "/users_pb.UserService/create", "/roles_pb.RoleService/getList"}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok || slices.Contains(withoutAuth, info.FullMethod) {
			return handler(ctx, req)
		}

		values := md["authorization"]
		if len(values) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing token")
		}

		token := values[0]
		claims, err := a.paseto.VerifyToken(token)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "token is invalid")
		}

		ctx = metadata.NewOutgoingContext(
			ctx, metadata.Pairs("userID", claims.Subject))

		return handler(ctx, req)
	}
}
