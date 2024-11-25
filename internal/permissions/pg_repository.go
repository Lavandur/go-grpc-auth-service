package permissions

import (
	"auth-service/internal/models"
	"context"
)

type PermissionRepository interface {
	GetRolePermissions(ctx context.Context, id string) ([]string, error)
	SetRolePermissions(ctx context.Context, id string, permission []string) (bool, error)

	GetPermissions(ctx context.Context) ([]*models.Permission, error)
	GetPermissionByID(ctx context.Context, id string) (*models.Permission, error)
	AddPermission(ctx context.Context, data *models.Permission) (*models.Permission, error)
	DeletePermission(ctx context.Context, id string) (bool, error)
}
