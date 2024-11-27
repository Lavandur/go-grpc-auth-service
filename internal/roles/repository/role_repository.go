package repository

import (
	"auth-service/internal/common"
	"auth-service/internal/models"
	"auth-service/internal/roles"
	"context"
	"errors"
	"github.com/doug-martin/goqu/v9"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type roleRepository struct {
	db     *pgxpool.Pool
	logger *logrus.Logger
}

func NewRoleRepository(db *pgxpool.Pool, logger *logrus.Logger) roles.RoleRepository {
	return &roleRepository{
		db:     db,
		logger: logger,
	}
}

func (r *roleRepository) GetByName(ctx context.Context, name string) (*models.Role, error) {
	r.logger.WithField("name", name).Debug("Get role by name")

	query, args, err := goqu.From("roles").
		Where(goqu.Ex{"name": name}).
		Limit(1).
		ToSQL()
	if err != nil {
		return nil, err
	}

	var role models.Role
	err = r.db.
		QueryRow(ctx, query, args...).
		Scan(&role.RoleID, &role.Name, &role.Description, &role.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}

	return &role, nil
}

func (r *roleRepository) GetByID(ctx context.Context, id string) (*models.Role, error) {
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

func (r *roleRepository) GetList(
	ctx context.Context,
	filter *models.RoleFilter,
	pagination *common.Pagination,
) ([]*models.Role, error) {

	roleList := make([]*models.Role, 0)

	whereList := r.getWhereList(filter)
	query, _, err := common.GetPagination(
		goqu.From("roles").Where(whereList...),
		pagination,
	).ToSQL()

	rows, err := r.db.Query(ctx, query)
	defer rows.Close()
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return roleList, nil
		}
		return nil, err
	}

	for rows.Next() {
		var role models.Role
		err = rows.Scan(
			&role.RoleID,
			&role.Name,
			&role.Description,
			&role.CreatedAt)
		if err != nil {
			return nil, err
		}
		roleList = append(roleList, &role)
	}
	return roleList, nil
}

func (r *roleRepository) Create(ctx context.Context, data *models.Role) (*models.Role, error) {
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

func (r *roleRepository) Update(ctx context.Context, data *models.Role) (*models.Role, error) {
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

func (r *roleRepository) Delete(ctx context.Context, id string) error {
	args := pgx.NamedArgs{
		"roleID": id,
	}

	err := r.db.QueryRow(ctx, deleteRole, args).Scan()
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return err
		}
	}
	return nil
}

func (r *roleRepository) getWhereList(filter *models.RoleFilter) []goqu.Expression {
	whereList := make([]goqu.Expression, 0)
	if filter == nil {
		return whereList
	}

	if filter.RoleID != nil {
		whereList = append(whereList, goqu.Ex{
			"id": filter.RoleID,
		})
	}
	if filter.Name != nil {
		whereList = append(whereList, goqu.Ex{
			"name": filter.Name,
		})
	}

	return whereList
}
