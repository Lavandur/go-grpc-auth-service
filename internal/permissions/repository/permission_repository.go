package repository

import (
	"auth-service/internal/models"
	"auth-service/internal/permissions"
	"context"
	"encoding/json"
	"errors"
	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
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

func (p *permissionRepository) GetPermissions(ctx context.Context) ([]*models.Permission, error) {
	p.logger.Debug("Getting permission-list", ctx)

	query, _, err := goqu.From("permissions").ToSQL()
	if err != nil {
		return nil, err
	}

	rows, err := p.db.Query(ctx, query)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	permissionList, err := pgx.CollectRows(
		rows, pgx.RowToStructByName[*models.Permission])
	if err != nil {
		return nil, err
	}

	return permissionList, nil
}

func (p *permissionRepository) AddPermission(ctx context.Context, data *models.PermissionInput) (*models.Permission, error) {
	p.logger.Debug("Adding permission", ctx)

	perm := models.Permission{
		PermissionID: uuid.New().String(),
		Name:         data.Name,
		Description:  data.Description,
	}

	description, err := json.Marshal(data.Description)
	if err != nil {
		return nil, err
	}
	query, _, err := goqu.Insert("permissions").Rows(
		goqu.Record{"id": perm.PermissionID, "name": perm.Name, "description": description},
	).ToSQL()
	if err != nil {
		return nil, err
	}

	err = p.db.QueryRow(ctx, query).Scan(
		&perm.PermissionID,
		&perm.Name,
		&perm.Description,
	)
	if err != nil {
		return nil, err
	}

	return &perm, nil
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
