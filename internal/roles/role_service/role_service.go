package role_service

import (
	"auth-service/internal/common"
	"auth-service/internal/models"
	"auth-service/internal/roles"
	"auth-service/pkg/config"
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type roleService struct {
	roleRepository roles.RoleRepository
	logger         *logrus.Logger

	defaultRoleName string
}

func NewRoleService(
	roleRepository roles.RoleRepository,
	logger *logrus.Logger,
	config *config.Config,
) roles.RoleService {

	return &roleService{
		roleRepository:  roleRepository,
		logger:          logger,
		defaultRoleName: config.DefaultRoleName,
	}
}

func (r *roleService) GetByID(ctx context.Context, id string) (*models.Role, error) {
	r.logger.Debugf("Get role by id: %s", id)

	role, err := r.roleRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return role, nil
}

func (r *roleService) GetList(ctx context.Context, filter *models.RoleFilter, pagination *common.Pagination) ([]*models.Role, error) {
	list, err := r.roleRepository.GetList(ctx, filter, pagination)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (r *roleService) Create(ctx context.Context, data *models.RoleInput) (*models.Role, error) {
	r.logger.Debugf("Create role with data: %+v", data)

	exists, err := r.roleRepository.GetByName(ctx, data.Title)
	if !(errors.Is(err, common.ErrNotFound) && exists == nil) {
		r.logger.Errorf("Failed creating role with name %s", data.Title)
		return exists, err
	}

	role, err := r.roleRepository.Create(
		ctx, &models.Role{
			RoleID:      uuid.NewString(),
			Title:       data.Title,
			Description: *data.Description,
			CreatedAt:   common.GetCurrentTime(),
		})
	if err != nil {
		return nil, err
	}

	return role, nil
}

func (r *roleService) Update(ctx context.Context, id string, data *models.RoleUpdateInput) (*models.Role, error) {
	r.logger.Debugf("Update role with data: %+v", data)

	old, err := r.roleRepository.GetByID(ctx, id)
	if old == nil || err != nil {
		r.logger.Errorf("Failed updating role with id %s", id)
		return old, err
	}

	updated, err := r.roleRepository.Update(ctx, data.ToUpdatedModel(old))
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (r *roleService) Delete(ctx context.Context, id string) error {
	r.logger.Debugf("Delete role by id: %s", id)

	if err := r.roleRepository.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}

func (r *roleService) GetDefaultRole(ctx context.Context) *models.Role {
	exists, err := r.roleRepository.GetByName(ctx, r.defaultRoleName)
	if !(errors.Is(err, common.ErrNotFound) && exists == nil) {
		r.logger.Debugf("Default role with name: %s already exists", r.defaultRoleName)
		return exists
	}

	// creating a new role if not exists
	role, err := r.roleRepository.Create(
		ctx, &models.Role{
			RoleID: uuid.NewString(),
			Title:  r.defaultRoleName,
			Description: common.LocalizedString{
				"en": "Default role",
				"ru": "Стандартная роль"},
			CreatedAt: common.GetCurrentTime(),
		})
	if err != nil {
		return nil
	}

	return role
}
