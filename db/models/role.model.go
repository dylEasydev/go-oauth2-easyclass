package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/dylEasydev/go-oauth2-easyclass/utils"
	"github.com/dylEasydev/go-oauth2-easyclass/validators"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ScopeJSON struct {
	ScopeName     string `json:"scopeName"`
	ScopeDescript string `json:"scopeDescript"`
}
type ScopeData struct {
	Data []ScopeJSON `json:"data"`
}

// structure des role d'utilisateur (admin , teacher , student ...)
type Role struct {
	ID           uuid.UUID `gorm:"primarykey;type:uuid;default:uuid_generate_v4()"`
	RoleName     string    `gorm:"column:rolename;not null;unique" validate:"required,rowallowed"`
	RoleDescript string    `gorm:"column:roledescript"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Scopes []*Scope `gorm:"many2many:authpermission;"`
}

func (Role) TableName() string {
	return "roles"
}

func (role *Role) BeforeSave(tx *gorm.DB) error {
	return validators.ValidateStruct(role)
}

func (role *Role) AddScope(tx *gorm.DB) (err error) {
	name := fmt.Sprintf("ressources/scope_%s", strings.ToLower(role.RoleName))
	data, err := utils.ReadJSON[ScopeData](name)

	if err != nil {
		return fmt.Errorf("lecture scopes pour role %s : %w", role.RoleName, err)
	}

	var sliceName []string = make([]string, len(data.Data))
	for i, elem := range data.Data {
		sliceName[i] = elem.ScopeName
	}

	var sopes []Scope
	if err = tx.Where("scopename IN ?", sliceName).Find(&sopes).Error; err != nil {
		return
	}

	if len(sopes) > 0 {
		if err = tx.Model(role).Association("Scopes").Replace(&sopes); err != nil {
			return fmt.Errorf("erreur remplacement scopes : %w", err)
		}
	}
	return nil
}

func (role *Role) AfterSave(tx *gorm.DB) error {
	if err := role.AddScope(tx); err != nil {
		fmt.Printf("Erreur lors de la mise à jour des scopes: %v\n", err)
	}
	return nil
}
