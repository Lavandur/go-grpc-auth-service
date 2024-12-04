package auth

import (
	"auth-service/internal/models"
	"context"
)

type AuthService interface {
	Login(ctx context.Context, login, password string) (*models.AuthResponse, error)
	Register(ctx context.Context, login, password string) (*models.AuthResponse, error)
	HasPermission(ctx context.Context, permission string) bool
}
