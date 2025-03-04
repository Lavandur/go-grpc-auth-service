// Code generated by MockGen. DO NOT EDIT.
// Source: internal/roles/role_service.go

// Package role_service is a generated GoMock package.
package roles_mock

import (
	common "auth-service/internal/common"
	models "auth-service/internal/models"
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockRoleService is a mock of RoleService interface.
type MockRoleService struct {
	ctrl     *gomock.Controller
	recorder *MockRoleServiceMockRecorder
}

// MockRoleServiceMockRecorder is the mock recorder for MockRoleService.
type MockRoleServiceMockRecorder struct {
	mock *MockRoleService
}

// NewMockRoleService creates a new mock instance.
func NewMockRoleService(ctrl *gomock.Controller) *MockRoleService {
	mock := &MockRoleService{ctrl: ctrl}
	mock.recorder = &MockRoleServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRoleService) EXPECT() *MockRoleServiceMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockRoleService) Create(ctx context.Context, data *models.RoleInput) (*models.Role, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, data)
	ret0, _ := ret[0].(*models.Role)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockRoleServiceMockRecorder) Create(ctx, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockRoleService)(nil).Create), ctx, data)
}

// Delete mocks base method.
func (m *MockRoleService) Delete(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockRoleServiceMockRecorder) Delete(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockRoleService)(nil).Delete), ctx, id)
}

// GetByID mocks base method.
func (m *MockRoleService) GetByID(ctx context.Context, id string) (*models.Role, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(*models.Role)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockRoleServiceMockRecorder) GetByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockRoleService)(nil).GetByID), ctx, id)
}

// GetDefaultRole mocks base method.
func (m *MockRoleService) GetDefaultRole(ctx context.Context) *models.Role {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDefaultRole", ctx)
	ret0, _ := ret[0].(*models.Role)
	return ret0
}

// GetDefaultRole indicates an expected call of GetDefaultRole.
func (mr *MockRoleServiceMockRecorder) GetDefaultRole(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDefaultRole", reflect.TypeOf((*MockRoleService)(nil).GetDefaultRole), ctx)
}

// GetList mocks base method.
func (m *MockRoleService) GetList(ctx context.Context, filter *models.RoleFilter, pagination *common.Pagination) ([]*models.Role, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetList", ctx, filter, pagination)
	ret0, _ := ret[0].([]*models.Role)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetList indicates an expected call of GetList.
func (mr *MockRoleServiceMockRecorder) GetList(ctx, filter, pagination interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetList", reflect.TypeOf((*MockRoleService)(nil).GetList), ctx, filter, pagination)
}

// Update mocks base method.
func (m *MockRoleService) Update(ctx context.Context, id string, data *models.RoleUpdateInput) (*models.Role, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, id, data)
	ret0, _ := ret[0].(*models.Role)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockRoleServiceMockRecorder) Update(ctx, id, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockRoleService)(nil).Update), ctx, id, data)
}
