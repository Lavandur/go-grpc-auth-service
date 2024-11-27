package auth

import (
	"auth-service/internal/common"
	"auth-service/internal/grpc/pb"
	"auth-service/internal/models"
	"auth-service/internal/service"
	"auth-service/internal/users"
	"context"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AuthServiceImpl struct {
	pb.UnimplementedAuthServiceServer

	authService service.AuthService
	userService users.UserService
	logger      *logrus.Logger
}

func NewAuthService(
	authService service.AuthService,
	userService users.UserService,
	logger *logrus.Logger,
) pb.AuthServiceServer {
	return &AuthServiceImpl{
		authService: authService,
		userService: userService,
		logger:      logger,
	}
}

func (s *AuthServiceImpl) Register(ctx context.Context, request *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	s.logger.Infof("Registering a new user with login: %v", request.Login)

	user, err := s.userService.Create(ctx, s.registerUserReqToModel(request))
	if err != nil {
		return nil, err
	}

	return s.userModelToRegReq(user), nil
}

func (s *AuthServiceImpl) Login(ctx context.Context, request *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	s.logger.Infof("Loging a new user with login: %v", request.Login)

	user, token, err := s.authService.Login(ctx, request.Login, request.Password)
	if err != nil {
		return nil, err
	}

	s.logger.Infof("Successfully logged in user: %v ||| %v", token, user)
	return &pb.LoginUserResponse{
		User:            s.userModelToProto(user),
		PublicToken:     token.PublicToken,
		PublicExpiresAt: timestamppb.New(token.PublicTokenExpiry),
	}, nil
}

func (s *AuthServiceImpl) GetByID(ctx context.Context, id *pb.ID) (*pb.User, error) {
	s.logger.Infof("Getting user by id: %s", id.Id)

	user, err := s.userService.GetByID(ctx, id.Id)
	if err != nil {
		return nil, err
	}

	return s.userModelToProto(user), nil
}

func (s *AuthServiceImpl) GetList(ctx context.Context, params *pb.UserListParams) (*pb.ArrayUser, error) {
	s.logger.Infof("Getting list of users by filter: %v and by pagination %v", params.Filter, params.Pagination)

	filter, pagination := s.userParamsToModel(params)
	listUsers, err := s.userService.GetList(ctx, filter, pagination)
	if err != nil {
		return nil, err
	}

	var result = make([]*pb.User, 0)
	for _, user := range listUsers {
		result = append(result, s.userModelToProto(user))
	}

	return &pb.ArrayUser{
		User: result,
	}, nil
}

func (s *AuthServiceImpl) RefreshAccessToken(ctx context.Context, request *pb.AccessTokenRequest) (*pb.AccessTokenResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *AuthServiceImpl) Me(ctx context.Context, empty *pb.Empty) (*pb.User, error) {
	//TODO implement me
	panic("implement me")
}

func (s *AuthServiceImpl) registerUserReqToModel(data *pb.RegisterUserRequest) *models.UserInput {
	return &models.UserInput{
		Login:     data.Login,
		Password:  data.Password,
		Firstname: data.Firstname,
		Lastname:  data.Lastname,
		Birthdate: data.Birthdate.AsTime(),
		Email:     data.Email,
		Gender:    data.Gender,
		RoleIDs:   data.RoleIDs,
	}
}

func (s *AuthServiceImpl) userModelToProto(data *models.User) *pb.User {
	return &pb.User{
		Login: data.Login,
		Person: &pb.Person{
			Firstname: data.Person.Firstname,
			Lastname:  data.Person.Lastname,
			Email:     data.Person.Email,
			Gender:    data.Person.Gender,
		},
	}
}

func (s *AuthServiceImpl) userModelToRegReq(data *models.User) *pb.RegisterUserResponse {
	return &pb.RegisterUserResponse{
		User: &pb.User{
			Login: data.Login,
			Person: &pb.Person{
				Firstname: data.Person.Firstname,
				Lastname:  data.Person.Lastname,
				Email:     data.Person.Email,
				Gender:    data.Person.Gender,
			},
		},
	}
}

func (s *AuthServiceImpl) userParamsToModel(data *pb.UserListParams) (*models.UserFilter, *common.Pagination) {
	if data == nil {
		return nil, nil
	}

	var newUserFilter *models.UserFilter
	if data.Filter.UserIDs != nil {
		newUserFilter.UserID = &data.Filter.UserIDs.UserID
	}
	if data.Filter.Logins != nil {
		newUserFilter.Login = &data.Filter.Logins.Login
	}
	if data.Filter.Emails != nil {
		newUserFilter.Email = &data.Filter.Emails.Email
	}

	var newPagination *common.Pagination
	if data.Pagination.Offset != nil {
		offset := uint(*data.Pagination.Offset)
		newPagination.Offset = &offset
	}
	if data.Pagination.Size != nil {
		limit := uint(*data.Pagination.Size)
		newPagination.Size = &limit
	}
	if data.Pagination.OrderBy != nil {
		newPagination.OrderBy = data.Pagination.OrderBy
	}
	return newUserFilter, newPagination
}
