package repository

import (
	"auth-service/internal/common"
	"auth-service/internal/models"
	"auth-service/internal/permissions"
	"context"
	"encoding/json"
	"errors"
	"github.com/doug-martin/goqu/v9"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type permissionRepository struct {
	db *pgxpool.Pool
}

func NewPermissionRepository(db *pgxpool.Pool) permissions.PermissionRepository {
	return &permissionRepository{
		db: db,
	}
}

func (p *permissionRepository) GetPermissionByID(ctx context.Context, id string) (*models.Permission, error) {

	query, _, err := goqu.
		From("permissions").
		Where(goqu.Ex{"id": id}).
		ToSQL()
	if err != nil {
		return nil, err
	}

	permission := &models.Permission{}
	err = p.db.QueryRow(ctx, query).
		Scan(&permission.PermissionID, &permission.Title, &permission.Description)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, common.ErrNotFound
		}
		return nil, err
	}

	return permission, nil
}

func (p *permissionRepository) GetPermissions(ctx context.Context) ([]*models.Permission, error) {
	query, _, err := goqu.From("permissions").ToSQL()
	if err != nil {
		return nil, err
	}

	rows, err := p.db.Query(ctx, query)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	var permissionList []*models.Permission
	for rows.Next() {
		permission := &models.Permission{}
		err = rows.Scan(&permission.PermissionID, &permission.Title, &permission.Description)

		if err != nil {
			return permissionList, err
		}

		permissionList = append(permissionList, permission)
	}

	return permissionList, nil
}

func (p *permissionRepository) AddPermission(ctx context.Context, data *models.Permission) (*models.Permission, error) {
	description, err := json.Marshal(data.Description)
	if err != nil {
		return nil, err
	}

	query, _, err := goqu.Insert("permissions").
		Rows(goqu.Record{
			"id":          data.PermissionID,
			"title":       data.Title,
			"description": description,
		},
		).
		Returning(
			"id",
			"title",
			"description").
		ToSQL()
	if err != nil {
		return nil, err
	}

	result := models.Permission{}
	err = p.db.QueryRow(ctx, query).Scan(
		&result.PermissionID,
		&result.Title,
		&result.Description,
	)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (p *permissionRepository) DeletePermission(ctx context.Context, id string) error {
	query, _, err := goqu.Delete("permissions").Where(goqu.Ex{"id": id}).ToSQL()
	if err != nil {
		return err
	}

	err = p.db.QueryRow(ctx, query).Scan()
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return err
		}
	}

	return nil
}

func (p *permissionRepository) GetRolePermissions(ctx context.Context, id string) ([]string, error) {
	args := pgx.NamedArgs{
		"roleID": id,
	}

	row := p.db.QueryRow(ctx, getPermissionsByID, args)

	return p.fetchPermissions(row)
}

func (p *permissionRepository) SetRolePermissions(ctx context.Context, id string, permissions []string) error {
	args := pgx.NamedArgs{
		"roleID":     id,
		"permission": permissions,
	}

	err := p.db.QueryRow(ctx, clearPermission, args).Scan()
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return err
	}

	err = p.db.QueryRow(ctx, setPermissions, args).Scan()
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return err
	}
	return nil
}

func (p *permissionRepository) fetchPermissions(row pgx.Row) ([]string, error) {
	var permissionList []string

	err := row.Scan(nil, &permissionList)
	if err != nil {
		return nil, err
	}

	return permissionList, nil
}
