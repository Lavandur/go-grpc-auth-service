package repository

import (
	"auth-service/internal/common"
	"auth-service/internal/models"
	"auth-service/internal/roles"
	"auth-service/internal/users"
	"context"
	"errors"
	"github.com/doug-martin/goqu/v9"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"time"
)

type usersRepository struct {
	db             *pgxpool.Pool
	roleRepository roles.RoleRepository
	logger         *logrus.Logger
}

func NewUsersRepository(pool *pgxpool.Pool, rolesRepository roles.RoleRepository, logger *logrus.Logger) users.UserRepository {
	return &usersRepository{
		db:             pool,
		roleRepository: rolesRepository,
		logger:         logger,
	}
}

func (r *usersRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	args := pgx.NamedArgs{
		"userID": id,
	}

	row := r.db.QueryRow(ctx, getUserByID, args)

	user, err := r.fetchUser(ctx, row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, common.ErrNotFound
		}
		return nil, err
	}
	return user, nil
}

func (r *usersRepository) GetList(ctx context.Context, filter *models.UserFilter) ([]*models.User, error) {

	usersList := make([]*models.User, 0)
	whereList, err := r.getWhereList(filter)
	if err != nil {
		return nil, err
	}

	query, _, _ := goqu.From("users").Where(whereList...).ToSQL()

	rows, err := r.db.Query(ctx, query)
	defer rows.Close()
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, common.ErrNotFound
		}
		return nil, err
	}

	for rows.Next() {
		user, err := r.fetchUser(ctx, rows)
		if err != nil {
			return nil, err
		}
		usersList = append(usersList, user)
	}

	return usersList, nil
}

func (r *usersRepository) Create(ctx context.Context, data *models.User) (*models.User, error) {
	var roleIDs []string
	for _, role := range data.Roles {
		roleIDs = append(roleIDs, role.RoleID)
	}

	args := pgx.NamedArgs{
		"userID":                data.UserID,
		"login":                 data.Login,
		"visibleID":             data.VisibleID,
		"hashedPassword":        data.HashedPassword,
		"person":                data.Person,
		"roles":                 roleIDs,
		"createdAt":             data.CreatedAt,
		"updatedAt":             data.UpdatedAt,
		"deletedAt":             data.DeletedAt,
		"lastPasswordRestoreAt": data.LastPasswordRestoreAt,
	}

	row := r.db.QueryRow(ctx, insertUser, args)

	return r.fetchUser(ctx, row)
}

func (r *usersRepository) Update(ctx context.Context, data *models.User) (*models.User, error) {
	roleIDs := make([]string, 0)
	for _, role := range data.Roles {
		roleIDs = append(roleIDs, role.RoleID)
	}

	args := pgx.NamedArgs{
		"userID":                data.UserID,
		"login":                 data.Login,
		"visibleID":             data.VisibleID,
		"hashedPassword":        data.HashedPassword,
		"person":                data.Person,
		"roles":                 roleIDs,
		"updatedAt":             data.UpdatedAt,
		"deletedAt":             data.DeletedAt,
		"lastPasswordRestoreAt": data.LastPasswordRestoreAt,
	}

	row := r.db.QueryRow(ctx, updateUser, args)

	return r.fetchUser(ctx, row)
}

func (r *usersRepository) Delete(ctx context.Context, id string) (*models.User, error) {
	args := pgx.NamedArgs{
		"userID":    id,
		"deletedAt": time.Now().UTC().Truncate(time.Microsecond),
	}

	row := r.db.QueryRow(ctx, deleteUser, args)

	return r.fetchUser(ctx, row)
}

func (r *usersRepository) getWhereList(filter *models.UserFilter) ([]goqu.Expression, error) {
	whereList := make([]goqu.Expression, 0)
	if filter == nil {
		return whereList, nil
	}

	//email is not a column, he is a part of person
	if filter.Email != nil {
		whereList = append(whereList, goqu.Ex{
			"email": filter.Email,
		})
	}
	if filter.Login != nil {
		whereList = append(whereList, goqu.Ex{
			"login": filter.Login,
		})
	}
	if filter.UserID != nil {
		whereList = append(whereList, goqu.Ex{
			"userID": filter.UserID,
		})
	}

	return whereList, nil
}

func (r *usersRepository) fetchUser(ctx context.Context, row pgx.Row) (*models.User, error) {
	roleIDs := make([]string, 0)
	var user models.User
	err := row.Scan(&user.UserID, &user.Login, &user.VisibleID,
		&user.HashedPassword, &user.Person, &roleIDs,
		&user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
		&user.LastPasswordRestoreAt, &user.SearchIndex)
	if err != nil {
		return nil, err
	}

	for _, roleID := range roleIDs {
		role, err := r.roleRepository.GetByID(ctx, roleID)
		if err != nil {
			return nil, err
		}
		user.Roles = append(user.Roles, role)
	}

	return &user, nil
}
