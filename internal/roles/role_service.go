package roles

import (
	"auth-service/internal/common"
	"auth-service/internal/models"
	"context"
)

type RoleService interface {
	GetByID(ctx context.Context, id string) (*models.Role, error)
	GetList(ctx context.Context, filter *models.RoleFilter, pagination *common.Pagination) ([]*models.Role, error)
	GetDefaultRole(ctx context.Context) *models.Role

	Create(ctx context.Context, data *models.RoleInput) (*models.Role, error)
	Update(ctx context.Context, id string, data *models.RoleUpdateInput) (*models.Role, error)
	Delete(ctx context.Context, id string) error
}
