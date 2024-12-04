package roles_pb

import (
	"auth-service/internal/grpc/pb"
	"auth-service/internal/models"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToProto(role *models.Role) *Role {
	return &Role{
		Id:          &pb.ID{Id: role.RoleID},
		Title:       role.Title,
		Description: &pb.LocalizedString{Value: role.Description},
		CreatedAt:   timestamppb.New(role.CreatedAt),
	}
}

func (f *RoleFilter) ToModel() *models.RoleFilter {
	if f == nil {
		return nil
	}

	var filter models.RoleFilter

	if f.Ids != nil {
		filter.RoleID = &f.Ids.Value
	}
	if f.Names != nil {
		filter.Title = &f.Names.Value
	}

	return &filter
}

func (r *RoleInput) ToModel() *models.RoleInput {
	var description map[string]string
	if r.Description != nil {
		description = r.Description.Value
	}

	return &models.RoleInput{
		Title:       r.Title,
		Description: description,
	}
}

func (r *RoleUpdatedInput) ToModel() *models.RoleUpdateInput {
	var description map[string]string
	if r.Description != nil {
		description = r.Description.Value
	}

	return &models.RoleUpdateInput{
		Title:       r.Title,
		Description: description,
	}
}

func (p *PermissionInput) ToModel() *models.PermissionInput {
	return &models.PermissionInput{
		Title:       p.Title,
		Description: p.Description,
	}
}

func PermissionToProto(permission *models.Permission) *Permission {
	return &Permission{
		Id:          &pb.ID{Id: permission.PermissionID},
		Title:       permission.Title,
		Description: permission.Description,
	}
}

func ListToProto(data []*models.Role) *RoleList {
	list := make([]*Role, 0, len(data))
	for _, role := range data {
		role := ToProto(role)
		list = append(list, role)
	}

	return &RoleList{
		Roles: list,
	}
}

func PermissionListToProto(data []*models.Permission) *PermissionList {
	list := make([]*Permission, 0, len(data))
	for _, permission := range data {
		list = append(list, PermissionToProto(permission))
	}

	return &PermissionList{
		List: list,
	}
}
