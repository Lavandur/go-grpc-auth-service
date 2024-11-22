package users

import (
	"auth-service/internal/models"
	"context"
)

type UserRepository interface {
	GetByID(ctx context.Context, id string) (*models.User, error)
	GetList(ctx context.Context, filter *models.UserFilter) ([]*models.User, error)

	Create(ctx context.Context, user *models.User) (*models.User, error)
	Update(ctx context.Context, user *models.User) (*models.User, error)
	Delete(ctx context.Context, id string) (*models.User, error)
}
