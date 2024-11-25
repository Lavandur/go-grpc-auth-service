package permission_service

import (
	"auth-service/internal/models"
	"auth-service/internal/permissions"
	"context"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type permissionService struct {
	permissionRepository permissions.PermissionRepository
	logger               *logrus.Logger
}

func NewPermissionService(
	permissionRepository permissions.PermissionRepository,
	logger *logrus.Logger,
) permissions.PermissionService {
	return &permissionService{
		permissionRepository: permissionRepository,
		logger:               logger,
	}
}

func (p *permissionService) GetByID(ctx context.Context, id string) (*models.Permission, error) {
	p.logger.WithField("id", id).Debug("getting permission by id")

	result, err := p.permissionRepository.GetPermissionByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (p *permissionService) GetList(ctx context.Context) ([]*models.Permission, error) {
	p.logger.Debug("getting list of permissions")

	permissionList, err := p.permissionRepository.GetPermissions(ctx)
	if err != nil {
		return nil, err
	}

	return permissionList, nil
}

func (p *permissionService) Create(ctx context.Context, data *models.PermissionInput) (*models.Permission, error) {
	p.logger.WithField("permission", data).Debug("creating new permission")

	perm := models.Permission{
		PermissionID: uuid.New().String(),
		Name:         data.Name,
		Description:  data.Description,
	}

	result, err := p.permissionRepository.AddPermission(ctx, &perm)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (p *permissionService) GetRolePermissions(ctx context.Context, id string) ([]string, error) {
	p.logger.WithField("id", id).Debug("getting role permissions")

	rolePermList, err := p.permissionRepository.GetRolePermissions(ctx, id)
	if err != nil {
		return nil, err
	}

	return rolePermList, nil
}

func (p *permissionService) SetRolePermissions(ctx context.Context, id string, permission []string) (bool, error) {
	p.logger.WithField("id", id).Debug("setting role permissions")

	success, err := p.permissionRepository.SetRolePermissions(ctx, id, permission)
	if err != nil {
		return false, err
	}

	return success, nil
}
