package repository

import (
	"auth-service/internal/users/models"
	"auth-service/internal/users/repository/roles"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"time"
)

type UsersRepository interface {
	GetByID(ctx context.Context, id string) (*models.User, error)

	Create(ctx context.Context, user *models.User) (*models.User, error)
	Update(ctx context.Context, user *models.User) (*models.User, error)
	Delete(ctx context.Context, id string) (*models.User, error)
}

type usersRepository struct {
	db       *pgxpool.Pool
	rolesRep roles.RoleRepository
	logger   *logrus.Logger
}

func NewUsersRepository(pool *pgxpool.Pool, rolesRepository roles.RoleRepository, logger *logrus.Logger) UsersRepository {
	return &usersRepository{
		db:       pool,
		rolesRep: rolesRepository,
		logger:   logger,
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
			return nil, nil
		}
		return nil, err
	}
	return user, nil
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

func (r *usersRepository) fetchUser(ctx context.Context, row pgx.Row) (*models.User, error) {
	var roleIDs []string
	var user models.User
	err := row.Scan(&user.UserID, &user.Login, &user.VisibleID,
		&user.HashedPassword, &user.Person, &roleIDs,
		&user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
		&user.LastPasswordRestoreAt, &user.SearchIndex)
	if err != nil {
		return nil, err
	}

	for _, roleID := range roleIDs {
		role, err := r.rolesRep.GetByID(ctx, roleID)
		if err != nil {
			return nil, err
		}
		user.Roles = append(user.Roles, role)
	}

	return &user, nil
}
