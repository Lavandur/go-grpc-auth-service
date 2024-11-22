package auth

import (
	"auth-service/internal/grpc/pb"
	"auth-service/internal/models"
	"auth-service/internal/users"
	"context"
	"github.com/sirupsen/logrus"
)

type AuthServiceImpl struct {
	pb.UnimplementedAuthServiceServer

	userService users.UserService
	paseto      *PasetoAuth
	logger      *logrus.Logger
}

func NewAuthService(
	userService users.UserService,
	paseto *PasetoAuth,
	logger *logrus.Logger,
) pb.AuthServiceServer {
	return &AuthServiceImpl{
		userService: userService,
		paseto:      paseto,
		logger:      logger,
	}
}

func (s *AuthServiceImpl) Register(ctx context.Context, request *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	s.logger.Infof("Registering a new user with login: %v", request.Login)

	user, err := s.userService.Create(ctx, s.registerUserReqToModel(request))
	if err != nil {
		return nil, err
	}

	return s.userModelToProto(user), nil
}

func (s *AuthServiceImpl) Login(ctx context.Context, request *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	s.logger.Infof("Loging a new user with login: %v", request.Login)

	data := s.loginUserReqToModel(request)
	user, err := s.userService.Login(ctx, data.Login, data.Password)
	if err != nil {
		return nil, err
	}

	s.logger.Infof("Access token are not implemented. Successfully logged in user: %v", user)
	return nil, nil
}

func (s *AuthServiceImpl) GetByID(ctx context.Context, id *pb.ID) (*pb.User, error) {
	s.logger.Infof("Getting user by id: %s", id.Id)

	user, err := s.userService.GetByID(ctx, id.Id)
	if err != nil {
		return nil, err
	}

	s.logger.WithFields(
		logrus.Fields{
			"Warning": "Casting user to proto aren't implemented.",
			"Result":  user,
		}).Info()

	return nil, nil
}

func (s *AuthServiceImpl) GetList(ctx context.Context, filter *pb.UserFilter) (*pb.ArrayUser, error) {
	s.logger.Infof("Getting list of users by filter: %v", filter)

	listUsers, err := s.userService.GetList(ctx, s.userFilterToModel(filter))
	if err != nil {
		return nil, err
	}

	s.logger.WithFields(
		logrus.Fields{
			"Warning": "Casting list of users to proto aren't implemented.",
			"Result":  listUsers,
		}).Info()

	return nil, nil
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

func (s *AuthServiceImpl) userModelToProto(data *models.User) *pb.RegisterUserResponse {
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

func (s *AuthServiceImpl) loginUserReqToModel(data *pb.LoginUserRequest) *models.UserLogin {
	return &models.UserLogin{
		Login:    data.Login,
		Password: data.Password,
	}
}

func (s *AuthServiceImpl) userFilterToModel(data *pb.UserFilter) *models.UserFilter {

	var userIds *[]string
	if data.UserIDs != nil {
		userIds = &data.UserIDs.UserID
	}
	var logins *[]string
	if data.Logins != nil {
		logins = &data.Logins.Login
	}
	var emails *[]string
	if data.Emails != nil {
		emails = &data.Emails.Email
	}

	return &models.UserFilter{
		UserID: userIds,
		Login:  logins,
		Email:  emails,
	}
}
