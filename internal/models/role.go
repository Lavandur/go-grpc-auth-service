package models

import (
	"auth-service/internal/common"
	"time"
)

type Permission struct {
	PermissionID string                 `json:"permissionID" db:"id"`
	Title        string                 `json:"title" db:"title"`
	Description  common.LocalizedString `json:"description" db:"description"`
}

type PermissionInput struct {
	Title       string
	Description common.LocalizedString
}

type Role struct {
	RoleID      string                 `json:"roleID"`
	Title       string                 `json:"title"`
	Description common.LocalizedString `json:"description,omitempty"`
	CreatedAt   time.Time              `json:"createdAt"`
}

type RoleInput struct {
	Title       string
	Description common.LocalizedString
}

type RoleUpdateInput struct {
	Title       *string
	Description common.LocalizedString
}

func (r *RoleUpdateInput) ToUpdatedModel(oldRole *Role) *Role {
	if r.Title != nil {
		oldRole.Title = *r.Title
	}
	if r.Description != nil {
		oldRole.Description = r.Description
	}

	return oldRole
}

type RoleFilter struct {
	RoleID *[]string
	Title  *[]string
}
