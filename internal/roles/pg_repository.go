package roles

import (
	"auth-service/internal/common"
	"auth-service/internal/models"
	"context"
)

type RoleRepository interface {
	GetByID(ctx context.Context, id string) (*models.Role, error)
	GetList(ctx context.Context, filter *models.RoleFilter, pagination *common.Pagination) ([]*models.Role, error)
	GetByName(ctx context.Context, name string) (*models.Role, error)

	Create(ctx context.Context, data *models.Role) (*models.Role, error)
	Update(ctx context.Context, data *models.Role) (*models.Role, error)
	Delete(ctx context.Context, id string) error
}
