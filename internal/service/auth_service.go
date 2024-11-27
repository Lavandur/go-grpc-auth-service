package service

import (
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

type AuthService interface {
	Login(ctx context.Context, login, password string) (*models.User, *models.AuthResponse, error)
}

type authService struct {
	paseto      PasetoAuth
	userService users.UserService

	logger *logrus.Logger
}

func NewAuthServiceImpl(
	paseto PasetoAuth,
	userService users.UserService,
	logger *logrus.Logger,
) AuthService {
	return &authService{
		paseto:      paseto,
		userService: userService,
		logger:      logger,
	}
}

func (a *authService) Login(ctx context.Context, login, password string) (*models.User, *models.AuthResponse, error) {
	a.logger.WithField("login", login).Debug("Login user")

	logins := []string{login}
	filter := &models.UserFilter{Login: &logins}
	list, err := a.userService.GetList(ctx, filter, nil)
	if err != nil || len(list) == 0 {
		a.logger.
			WithField("login", login).
			Error("user by login not found")
		return nil, nil, common.ErrNotFound
	}

	user := list[0]

	if !a.checkPasswordHash(password, user.HashedPassword) {
		a.logger.WithFields(
			logrus.Fields{
				"userID":   user.UserID,
				"login":    login,
				"password": password,
			},
		).Error("failed comparing passwords")
		return nil, nil, ErrWrongPassword
	}

	data := TokenData{
		Subject:  user.UserID,
		Duration: 60 * time.Second,
		AdditionalClaims: AdditionalClaims{
			ID:   user.UserID,
			Role: user.Login,
		},
		Footer: Footer{},
	}

	token, err := a.paseto.NewToken(data)
	if err != nil {
		return nil, nil, err
	}
	claims, err := a.paseto.VerifyToken(token)
	if err != nil {
		return nil, nil, err
	}

	response := &models.AuthResponse{
		PublicToken:       token,
		PublicTokenExpiry: claims.Expiration.UTC(),
	}

	return user, response, nil
}

func (a *authService) checkPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}
