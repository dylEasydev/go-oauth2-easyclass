package models

import (
	"time"

	"github.com/dylEasydev/go-oauth2-easyclass/validators"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// strcutures des permission d'utilisateur en fonction des role
type Scope struct {
	ID            uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	ScopeName     string    `gorm:"column:scope_name;not null;uniqueIndex" validate:"required,name"`
	ScopeDescript string    `gorm:"column:scope_descript;"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Roles []*Role `gorm:"many2many:authpermission;"`
}

func (Scope) TableName() string {
	return "scopes"
}

func (scope *Scope) BeforeSave(tx *gorm.DB) error {
	return validators.ValidateStruct(scope)
}
