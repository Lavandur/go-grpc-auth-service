package roles

import (
	"auth-service/internal/common"
	"auth-service/internal/users/models"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type RoleRepository interface {
	GetByID(ctx context.Context, id string) (*models.Role, error)
	GetList(ctx context.Context, pagination common.Pagination) ([]*models.Role, error)

	Create(ctx context.Context, data *models.Role) (*models.Role, error)
	Update(ctx context.Context, data *models.Role) (*models.Role, error)
	Delete(ctx context.Context, id string) error
}

type roleRepository struct {
	db     *pgxpool.Pool
	logger *logrus.Logger
}

func NewRoleRepository(db *pgxpool.Pool, logger *logrus.Logger) RoleRepository {
	return &roleRepository{
		db:     db,
		logger: logger,
	}
}

func (r roleRepository) GetByID(ctx context.Context, id string) (*models.Role, error) {
	args := pgx.NamedArgs{
		"roleID": id,
	}

	var role models.Role
	err := r.db.
		QueryRow(ctx, getRoleByID, args).
		Scan(&role.RoleID, &role.Name, &role.Description, &role.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r roleRepository) GetList(ctx context.Context, pagination common.Pagination) ([]*models.Role, error) {
	args := pgx.NamedArgs{
		"order_by": pagination.OrderBy,
		"offset":   pagination.Offset,
		"limit":    pagination.Size,
	}

	rows, err := r.db.Query(ctx, getRoles, args)
	defer rows.Close()
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	var roles []*models.Role
	for rows.Next() {
		var role models.Role
		err = rows.Scan(&role.RoleID, &role.Name, &role.Description, &role.CreatedAt)
		if err != nil {
			return nil, err
		}
		roles = append(roles, &role)
	}
	return roles, nil
}

func (r roleRepository) Create(ctx context.Context, data *models.Role) (*models.Role, error) {
	args := pgx.NamedArgs{
		"roleID":      data.RoleID,
		"name":        data.Name,
		"description": data.Description,
		"createdAt":   data.CreatedAt,
	}

	var role models.Role
	err := r.db.
		QueryRow(ctx, createRole, args).
		Scan(&role.RoleID, &role.Name, &role.Description, &role.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r roleRepository) Update(ctx context.Context, data *models.Role) (*models.Role, error) {
	args := pgx.NamedArgs{
		"roleID":      data.RoleID,
		"name":        data.Name,
		"description": data.Description,
		"createdAt":   data.CreatedAt,
	}

	var role models.Role
	err := r.db.
		QueryRow(ctx, updateRole, args).
		Scan(&role.RoleID, &role.Name, &role.Description, &role.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r roleRepository) Delete(ctx context.Context, id string) error {
	args := pgx.NamedArgs{
		"roleID": id,
	}

	err := r.db.QueryRow(ctx, deleteRole, args).Scan()
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		return err
	}
	return nil
}
