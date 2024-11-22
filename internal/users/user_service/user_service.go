package user_service

import (
	"auth-service/internal/models"
	"auth-service/internal/roles"
	"auth-service/internal/users"
	"context"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type userService struct {
	userRepository users.UserRepository
	roleRepository roles.RoleRepository

	logger *logrus.Logger
}

func NewUserService(
	userRepository users.UserRepository,
	roleRepository roles.RoleRepository,
	logger *logrus.Logger,
) users.UserService {
	return &userService{
		userRepository: userRepository,
		roleRepository: roleRepository,
		logger:         logger,
	}
}

func (u *userService) GetByID(ctx context.Context, id string) (*models.User, error) {
	u.logger.Infof("Getting user by ID: %s", id)

	user, err := u.userRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userService) GetList(ctx context.Context, filter *models.UserFilter) ([]*models.User, error) {
	u.logger.Infof("Getting users by filter: %+v", filter)

	listUsers, err := u.userRepository.GetList(ctx, filter)
	if err != nil {
		return nil, err
	}

	return listUsers, nil
}

func (u *userService) Create(ctx context.Context, data *models.UserInput) (*models.User, error) {
	uid, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	userRoles := make([]*models.Role, 0)
	for _, id := range data.RoleIDs {
		role, _ := u.roleRepository.GetByID(ctx, id)
		/*if err != nil {
			return nil, err
		}*/
		userRoles = append(userRoles, role)
	}

	user := &models.User{
		UserID:         uid.String(),
		Login:          data.Login,
		VisibleID:      data.Login,
		HashedPassword: u.hashPassword(data.Password),
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

func (u *userService) Login(ctx context.Context, login, password string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userService) Update(ctx context.Context, id string, data *models.UserUpdateInput) (*models.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userService) Delete(ctx context.Context, id string) (*models.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userService) hashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 6)

	return string(bytes)
}

func (u *userService) checkPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}
