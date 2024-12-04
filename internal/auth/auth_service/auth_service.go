package auth_service

import (
	"auth-service/internal/auth"
	"auth-service/internal/models"
	"auth-service/internal/permissions"
	"auth-service/internal/users"
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"slices"
)

var (
	ErrWrongPassword   = errors.New("wrong password")
	ErrNotEnoughRights = errors.New("not enough rights")
)

type authService struct {
	paseto            PasetoAuth
	userService       users.UserService
	permissionService permissions.PermissionService

	logger *logrus.Logger
}

func NewAuthServiceImpl(
	paseto PasetoAuth,
	userService users.UserService,
	permissionService permissions.PermissionService,
	logger *logrus.Logger,
) auth.AuthService {
	return &authService{
		paseto:            paseto,
		userService:       userService,
		permissionService: permissionService,
		logger:            logger,
	}
}

func (a *authService) Register(ctx context.Context, login, password string) (*models.AuthResponse, error) {
	a.logger.Debugf("Register user with login: %s", login)

	/*data

	token, claims, err := a.paseto.NewToken()
	if err != nil {
		return nil, err
	}*/
	return nil, nil
}

func (a *authService) Login(ctx context.Context, login, password string) (*models.AuthResponse, error) {
	a.logger.WithField("login", login).Debug("Login user")

	user, err := a.userService.GetByLogin(ctx, login)
	if err != nil {
		return nil, err
	}

	if !a.checkPasswordHash(password, user.HashedPassword) {
		a.logger.WithFields(
			logrus.Fields{
				"userID":   user.UserID,
				"login":    login,
				"password": password,
			},
		).Error(ErrWrongPassword)
		return nil, ErrWrongPassword
	}

	token, claims, err := a.paseto.NewToken(user)
	if err != nil {
		return nil, err
	}

	response := &models.AuthResponse{
		PublicToken:       token,
		PublicTokenExpiry: claims.Expiration.UTC(),
	}

	return response, nil
}

func (a *authService) HasPermission(ctx context.Context, permission string) bool {
	userID, ok := ctx.Value("userID").(string)
	if !ok || userID == "" {
		a.logger.Debugf("userID not found or userID is blank in ctx")
		return false
	}

	a.logger.Debugf("Check permission: %s for user with id: %s", permission, userID)

	user, err := a.userService.GetByID(ctx, userID)
	if err != nil {
		return false
	}

	for _, role := range user.Roles {
		permList, err := a.permissionService.GetRolePermissions(ctx, role.RoleID)
		if err != nil {
			return false
		}

		if slices.Contains(permList, permission) {
			return true
		}
	}

	return false
}

func (a *authService) checkPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}
