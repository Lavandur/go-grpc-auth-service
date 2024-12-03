package delivery

import (
	"auth-service/internal/grpc/pb"
	"auth-service/internal/grpc/pb/roles_pb"
	"auth-service/internal/permissions"
	"auth-service/internal/roles"
	"context"
	"github.com/sirupsen/logrus"
)

type RoleGRPCService struct {
	roles_pb.UnimplementedRoleServiceServer

	roleService       roles.RoleService
	permissionService permissions.PermissionService
	logger            *logrus.Logger
}

func NewRoleGRPC(
	roleService roles.RoleService,
	permissionService permissions.PermissionService,
	logger *logrus.Logger,
) *RoleGRPCService {
	return &RoleGRPCService{
		roleService:       roleService,
		permissionService: permissionService,
		logger:            logger,
	}
}

func (r *RoleGRPCService) GetByID(ctx context.Context, id *pb.ID) (*roles_pb.Role, error) {
	r.logger.Debugf("Get role by id: %s", id.GetId())

	role, err := r.roleService.GetByID(ctx, id.GetId())
	if err != nil {
		return nil, err
	}

	return roles_pb.ToProto(role), nil
}

func (r *RoleGRPCService) GetList(ctx context.Context, params *roles_pb.RoleListParams) (*roles_pb.RoleList, error) {
	r.logger.Debugf("Get role list with filter: %v pagination: %v", params.Filter, params.Pagination)

	filter := params.Filter.ToModel()
	pagination := params.Pagination.ToModel()
	list, err := r.roleService.GetList(ctx, filter, pagination)
	if err != nil {
		return nil, err
	}

	return roles_pb.ListToProto(list), nil
}

func (r *RoleGRPCService) Create(ctx context.Context, request *roles_pb.RoleCreateRequest) (*roles_pb.Role, error) {
	r.logger.Debugf("Create role request: %v", request)

	role, err := r.roleService.Create(ctx, request.Data.ToModel())
	if err != nil {
		return nil, err
	}

	return roles_pb.ToProto(role), nil
}

func (r *RoleGRPCService) Update(ctx context.Context, request *roles_pb.RoleUpdateRequest) (*roles_pb.Role, error) {
	r.logger.Debugf("Update role request: %v", request)

	role, err := r.roleService.Update(ctx, request.Id.GetId(), request.Data.ToModel())
	if err != nil {
		return nil, err
	}

	return roles_pb.ToProto(role), nil
}

func (r *RoleGRPCService) Delete(ctx context.Context, id *pb.ID) (*pb.IsSuccess, error) {
	r.logger.Debugf("Delete role by id: %s", id.GetId())

	err := r.roleService.Delete(ctx, id.GetId())
	if err != nil {
		return nil, err
	}

	return &pb.IsSuccess{Value: true}, nil
}

func (r *RoleGRPCService) GetRolePermissions(ctx context.Context, id *pb.ID) (*pb.ArrayString, error) {
	r.logger.Debugf("Get role permissions by role id: %s", id.GetId())

	rolePermissions, err := r.permissionService.GetRolePermissions(ctx, id.GetId())
	if err != nil {
		return nil, err
	}

	return &pb.ArrayString{Value: rolePermissions}, nil
}

func (r *RoleGRPCService) SetRolePermissions(ctx context.Context, request *roles_pb.SetRolePermissionsRequest) (*pb.IsSuccess, error) {
	r.logger.Debugf("Set role permissions with role id: %s", request.GetId())

	_, err := r.permissionService.SetRolePermissions(ctx,
		request.Id.GetId(),
		request.Permissions.GetValue())
	if err != nil {
		return nil, err
	}

	return &pb.IsSuccess{Value: true}, nil
}

func (r *RoleGRPCService) GetPermissionByID(ctx context.Context, id *pb.ID) (*roles_pb.Permission, error) {
	r.logger.Debugf("Get permission by id: %s", id.GetId())

	permission, err := r.permissionService.GetByID(ctx, id.GetId())
	if err != nil {
		return nil, err
	}

	return roles_pb.PermissionToProto(permission), nil
}

func (r *RoleGRPCService) GetPermissionList(ctx context.Context, _ *pb.Empty) (*roles_pb.PermissionList, error) {
	r.logger.Debugf("Get all permissions")

	list, err := r.permissionService.GetList(ctx)
	if err != nil {
		return nil, err
	}

	return roles_pb.PermissionListToProto(list), nil
}

func (r *RoleGRPCService) CreatePermission(ctx context.Context, request *roles_pb.PermissionCreateRequest) (*roles_pb.Permission, error) {
	r.logger.Debugf("Create permission request: %v", request.Data)

	permission, err := r.permissionService.Create(ctx, request.Data.ToModel())
	if err != nil {
		return nil, err
	}

	return roles_pb.PermissionToProto(permission), nil
}
