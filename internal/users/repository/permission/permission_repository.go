package permission

import (
	"auth-service/internal/users/models"
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type PermissionRepository interface {
	GetPermissionsByID(ctx context.Context, id string) ([]*models.Permission, error)
	SetPermissionByID(ctx context.Context, roleID, permission string) (bool, error)
}

type permissionRepository struct {
	db     *pgxpool.Pool
	logger *logrus.Logger
}

func NewPermissionRepository(db *pgxpool.Pool, logger *logrus.Logger) PermissionRepository {
	return &permissionRepository{
		db:     db,
		logger: logger,
	}
}

func (p permissionRepository) GetPermissionsByID(ctx context.Context, id string) ([]*models.Permission, error) {
	args := pgx.NamedArgs{
		"roleID": id,
	}

	var permissions []*models.Permission
	rows, err := p.db.Query(ctx, getPermissionByID, args)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var permission models.Permission
		err := rows.Scan(&permission.PermissionID, &permission.Name, &permission.Description)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, &permission)
	}

	return permissions, nil
}

func (p permissionRepository) SetPermissionByID(ctx context.Context, roleID, permission string) (bool, error) {
	args := pgx.NamedArgs{
		"roleID":     roleID,
		"permission": permission,
	}

	err := p.db.
		QueryRow(ctx, setPermission, args).
		Scan()
	if err != nil {
		return false, err
	}
	return true, nil
}
