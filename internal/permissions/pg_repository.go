package permissions

import "context"

type PermissionRepository interface {
	GetRolePermissions(ctx context.Context, id string) ([]string, error)
	SetRolePermissions(ctx context.Context, id string, permission []string) (bool, error)
}
