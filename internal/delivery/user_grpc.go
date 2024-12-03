package delivery

import (
	"auth-service/internal/common"
	"auth-service/internal/grpc/pb"
	"auth-service/internal/grpc/pb/users_pb"
	"auth-service/internal/users"
	"context"
	"errors"
	"github.com/sirupsen/logrus"
)

type UserGrpcService struct {
	users_pb.UnimplementedUserServiceServer

	userService users.UserService
	logger      *logrus.Logger
}

func NewUserGrpcService(
	userService users.UserService,
	logger *logrus.Logger,
) *UserGrpcService {
	return &UserGrpcService{
		userService: userService,
		logger:      logger,
	}
}

func (u *UserGrpcService) GetList(ctx context.Context, params *users_pb.UserListParams) (*users_pb.ArrayUser, error) {
	u.logger.Infof(
		"Get list of users by filter: %v and by pagination %v",
		params.Filter, params.Pagination)

	filter := params.Filter.ToModel()
	pagination := params.Pagination.ToModel()
	listUsers, err := u.userService.GetList(ctx, filter, pagination)
	if err != nil {
		return nil, err
	}

	var result = make([]*users_pb.User, 0)
	for _, user := range listUsers {
		result = append(result, users_pb.ToProto(user))
	}

	return &users_pb.ArrayUser{User: result}, nil
}

func (u *UserGrpcService) GetByID(ctx context.Context, id *pb.ID) (*users_pb.User, error) {
	u.logger.Infof("Get user by id: %s", id.Id)

	user, err := u.userService.GetByID(ctx, id.Id)
	if err != nil {
		return nil, err
	}

	return users_pb.ToProto(user), nil
}

func (u *UserGrpcService) Create(ctx context.Context, request *users_pb.CreateUserRequest) (*users_pb.User, error) {
	u.logger.Debugf("Create a new user with login: %v", request.Login)

	usr, err := u.userService.GetByLogin(ctx, request.Login)
	if !errors.Is(err, common.ErrNotFound) {
		return nil, err
	}
	if usr != nil {
		return nil, common.ErrLoginExists
	}

	user, err := u.userService.Create(ctx, request.ToModel())
	if err != nil {
		return nil, err
	}

	return users_pb.ToProto(user), nil
}

func (u *UserGrpcService) Update(ctx context.Context, request *users_pb.UpdateUserRequest) (*users_pb.User, error) {
	u.logger.Debugf("Update a user with id: %s", request.Id.Id)

	id, data := request.ToModel()
	user, err := u.userService.Update(ctx, id, data)
	if err != nil {
		return nil, err
	}

	return users_pb.ToProto(user), nil
}

func (u *UserGrpcService) DeleteByID(ctx context.Context, id *pb.ID) (*users_pb.User, error) {
	u.logger.Debugf("Delete a user with id: %s", id)

	user, err := u.userService.Delete(ctx, id.Id)
	if err != nil {
		return nil, err
	}

	return users_pb.ToProto(user), nil
}

func (u *UserGrpcService) Me(ctx context.Context, empty *pb.Empty) (*users_pb.User, error) {
	//TODO implement me
	panic("implement me")
}
