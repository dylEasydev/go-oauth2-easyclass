package models

import "github.com/google/uuid"

type AuthPermission struct {
	RoleID  uuid.UUID `gorm:"primaryKey;uniqueIndex:idx_role_scope"`
	ScopeID uuid.UUID `gorm:"primaryKey;uniqueIndex:idx_role_scope"`
}

func (AuthPermission) TableName() string {
	return "authpermission"
}
