package repository

import (
	"auth-service/internal/permissions"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type permissionRepository struct {
	db     *pgxpool.Pool
	logger *logrus.Logger
}

func NewPermissionRepository(db *pgxpool.Pool, logger *logrus.Logger) permissions.PermissionRepository {
	return &permissionRepository{
		db:     db,
		logger: logger,
	}
}

func (p *permissionRepository) GetRolePermissions(ctx context.Context, id string) ([]string, error) {
	args := pgx.NamedArgs{
		"roleID": id,
	}

	row := p.db.QueryRow(ctx, getPermissionsByID, args)

	return p.fetchPermissions(row)
}

func (p *permissionRepository) SetRolePermissions(ctx context.Context, id string, permissions []string) (bool, error) {
	args := pgx.NamedArgs{
		"roleID":     id,
		"permission": permissions,
	}

	err := p.db.QueryRow(ctx, clearPermission, args).Scan()
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return false, err
	}

	err = p.db.QueryRow(ctx, setPermissions, args).Scan()
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return false, err
	}
	return true, nil
}

func (p *permissionRepository) fetchPermissions(row pgx.Row) ([]string, error) {
	var permissionList []string

	err := row.Scan(nil, &permissionList)
	if err != nil {
		return nil, err
	}

	return permissionList, nil
}
