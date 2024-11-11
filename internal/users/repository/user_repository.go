package repository

import (
	"auth-service/internal/users/models"
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
	db     *pgxpool.Pool
	logger *logrus.Logger
}

func NewUsersRepository(pool *pgxpool.Pool, logger *logrus.Logger) UsersRepository {
	return &usersRepository{
		db:     pool,
		logger: logger,
	}
}

func (r *usersRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	args := pgx.NamedArgs{
		"userID": id,
	}

	var user models.User
	err := r.db.QueryRow(ctx, getUserByID, args).
		Scan(&user.UserID, &user.Login, &user.VisibleID,
			&user.HashedPassword, &user.Person, &user.Roles,
			&user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
			&user.LastPasswordRestoreAt, &user.SearchIndex)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *usersRepository) Create(ctx context.Context, data *models.User) (*models.User, error) {
	args := pgx.NamedArgs{
		"userID":                data.UserID,
		"login":                 data.Login,
		"visibleID":             data.VisibleID,
		"hashedPassword":        data.HashedPassword,
		"person":                data.Person,
		"roles":                 data.Roles,
		"createdAt":             data.CreatedAt,
		"updatedAt":             data.UpdatedAt,
		"deletedAt":             data.DeletedAt,
		"lastPasswordRestoreAt": data.LastPasswordRestoreAt,
	}

	var user models.User
	err := r.db.
		QueryRow(ctx, insertUser, args).
		Scan(&user.UserID, &user.Login, &user.VisibleID,
			&user.HashedPassword, &user.Person, &user.Roles,
			&user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
			&user.LastPasswordRestoreAt, &user.SearchIndex)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *usersRepository) Update(ctx context.Context, data *models.User) (*models.User, error) {
	args := pgx.NamedArgs{
		"userID":                data.UserID,
		"login":                 data.Login,
		"visibleID":             data.VisibleID,
		"hashedPassword":        data.HashedPassword,
		"person":                data.Person,
		"roles":                 data.Roles,
		"updatedAt":             data.UpdatedAt,
		"deletedAt":             data.DeletedAt,
		"lastPasswordRestoreAt": data.LastPasswordRestoreAt,
	}

	var user models.User
	err := r.db.
		QueryRow(ctx, updateUser, args).
		Scan(&user.UserID, &user.Login, &user.VisibleID,
			&user.HashedPassword, &user.Person, &user.Roles, &user.CreatedAt,
			&user.UpdatedAt, &user.DeletedAt, &user.LastPasswordRestoreAt,
			&user.SearchIndex)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *usersRepository) Delete(ctx context.Context, id string) (*models.User, error) {
	args := pgx.NamedArgs{
		"userID":    id,
		"deletedAt": time.Now(),
	}

	var user models.User
	err := r.db.QueryRow(ctx, deleteUser, args).
		Scan(&user.UserID, &user.Login, &user.VisibleID,
			&user.HashedPassword, &user.Person, &user.Roles, &user.CreatedAt,
			&user.UpdatedAt, &user.DeletedAt, &user.LastPasswordRestoreAt,
			&user.SearchIndex)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
