package user_service

import (
	"auth-service/internal/common"
	"auth-service/internal/models"
	"auth-service/internal/roles"
	"auth-service/internal/users"
	"context"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"time"
)

type userService struct {
	userRepository users.UserRepository
	roleService    roles.RoleService

	logger *logrus.Logger
}

func NewUserService(
	userRepository users.UserRepository,
	roleService roles.RoleService,
	logger *logrus.Logger,
) users.UserService {
	return &userService{
		userRepository: userRepository,
		roleService:    roleService,
		logger:         logger,
	}
}

func (u *userService) GetByID(ctx context.Context, id string) (*models.User, error) {
	u.logger.Infof("Get user by ID: %s", id)

	user, err := u.userRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userService) GetByLogin(ctx context.Context, login string) (*models.User, error) {
	u.logger.Infof("Get user by login: %s", login)

	user, err := u.userRepository.GetByLogin(ctx, login)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userService) GetList(ctx context.Context, filter *models.UserFilter, pagination *common.Pagination) ([]*models.User, error) {
	u.logger.Infof("Get users by filter: %+v", filter)

	listUsers, err := u.userRepository.GetList(ctx, filter, pagination)
	if err != nil {
		return nil, err
	}

	return listUsers, nil
}

func (u *userService) Create(ctx context.Context, data *models.UserInput) (*models.User, error) {

	userID, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	roleFilter := &models.RoleFilter{
		RoleID: &data.RoleIDs,
	}
	userRoles, err := u.roleService.GetList(ctx, roleFilter, nil)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		UserID:         userID.String(),
		Login:          data.Login,
		VisibleID:      data.Login,
		HashedPassword: common.HashPassword(data.Password),
		Person: models.Person{
			Firstname: data.Firstname,
			Lastname:  data.Lastname,
			Birthdate: data.Birthdate,
			Email:     data.Email,
			Gender:    data.Gender,
		},
		Roles:                 userRoles,
		CreatedAt:             time.Now().UTC().Truncate(time.Millisecond),
		UpdatedAt:             time.Now().UTC().Truncate(time.Millisecond),
		DeletedAt:             nil,
		LastPasswordRestoreAt: nil,
		SearchIndex:           nil,
	}
	res, err := u.userRepository.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (u *userService) Update(ctx context.Context, id string, data *models.UserUpdateInput) (*models.User, error) {
	u.logger.Debugf("Update user by ID: %s", id)

	user, err := u.userRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	data.ToUpdatedModel(user)

	filter := &models.RoleFilter{
		RoleID: &data.RoleIDs,
	}

	roleList, err := u.roleService.GetList(ctx, filter, nil)
	user.Roles = roleList

	result, err := u.userRepository.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u *userService) Delete(ctx context.Context, id string) (*models.User, error) {
	u.logger.Debugf("Delete user by ID: %s", id)

	user, err := u.userRepository.Delete(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}
