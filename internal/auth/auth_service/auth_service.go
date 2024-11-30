package auth_service

import (
	"auth-service/internal/auth"
	"auth-service/internal/common"
	"auth-service/internal/models"
	"auth-service/internal/users"
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var (
	ErrWrongPassword = errors.New("wrong password")
)

type authService struct {
	paseto      PasetoAuth
	userService users.UserService

	logger *logrus.Logger
}

func NewAuthServiceImpl(
	paseto PasetoAuth,
	userService users.UserService,
	logger *logrus.Logger,
) auth.AuthService {
	return &authService{
		paseto:      paseto,
		userService: userService,
		logger:      logger,
	}
}

func (a *authService) Register(ctx context.Context, login, password string) (*models.AuthResponse, error) {
	a.logger.Debugf("Register user with login: %s", login)

	/*data

	token, claims, err := a.paseto.NewToken()
	if err != nil {
		return nil, err
	}

	return*/
	return nil, nil
}

func (a *authService) Login(ctx context.Context, login, password string) (*models.AuthResponse, error) {
	a.logger.WithField("login", login).Debug("Login user")

	user, err := a.userService.GetByLogin(ctx, login)
	if err != nil {
		a.logger.
			WithField("login", login).
			Error("user by login not found")
		return nil, common.ErrNotFound
	}

	if !a.checkPasswordHash(password, user.HashedPassword) {
		a.logger.WithFields(
			logrus.Fields{
				"userID":   user.UserID,
				"login":    login,
				"password": password,
			},
		).Error("failed comparing passwords")
		return nil, ErrWrongPassword
	}

	data := models.TokenData{
		Subject:  user.UserID,
		Duration: 60 * time.Second,
		AdditionalClaims: models.AdditionalClaims{
			ID:   user.UserID,
			Role: user.Login,
		},
		Footer: models.Footer{},
	}

	token, claims, err := a.paseto.NewToken(data)
	if err != nil {
		return nil, err
	}

	response := &models.AuthResponse{
		PublicToken:       token,
		PublicTokenExpiry: claims.Expiration.UTC(),
	}

	return response, nil
}

func (a *authService) checkPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}
