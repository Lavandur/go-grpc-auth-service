package permissions

import (
	"auth-service/internal/models"
	"context"
)

type PermissionService interface {
	GetRolePermissions(ctx context.Context, id string) ([]string, error)
	SetRolePermissions(ctx context.Context, id string, permission []string) (bool, error)

	GetByID(ctx context.Context, id string) (*models.Permission, error)
	GetList(ctx context.Context) ([]*models.Permission, error)

	Create(ctx context.Context, data *models.PermissionInput) (*models.Permission, error)
}
