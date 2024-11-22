package models

import (
	"auth-service/internal/common"
	"time"
)

type Permission struct {
	PermissionID string                 `json:"permissionID" db:"id"`
	Name         string                 `json:"name" db:"name"`
	Description  common.LocalizedString `json:"description" db:"description"`
}

type PermissionInput struct {
	Name        string
	Description common.LocalizedString
}

type Role struct {
	RoleID      string                 `json:"roleID"`
	Name        string                 `json:"name"`
	Description common.LocalizedString `json:"description,omitempty"`
	CreatedAt   time.Time              `json:"createdAt"`
}
