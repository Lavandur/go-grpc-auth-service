package models

import (
	"auth-service/internal/common"
	"time"
)

type Permission struct {
	PermissionID string                 `json:"roleID"`
	Name         string                 `json:"roleName"`
	Description  common.LocalizedString `json:"description"`
}

type Role struct {
	RoleID      string                 `json:"roleID"`
	Name        string                 `json:"name"`
	Description common.LocalizedString `json:"description,omitempty"`
	CreatedAt   time.Time              `json:"createdAt"`
}
