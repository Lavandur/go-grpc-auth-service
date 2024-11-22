package users

import (
	"auth-service/internal/models"
	"context"
)

type UserService interface {
	GetByID(ctx context.Context, id string) (*models.User, error)
	GetList(ctx context.Context, filter *models.UserFilter) ([]*models.User, error)

	Create(ctx context.Context, data *models.UserInput) (*models.User, error)
	Update(ctx context.Context, id string, data *models.UserUpdateInput) (*models.User, error)
	Delete(ctx context.Context, id string) (*models.User, error)
	Login(ctx context.Context, login, password string) (bool, error)
}
