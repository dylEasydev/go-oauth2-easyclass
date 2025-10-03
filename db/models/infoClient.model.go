package models

import (
	"time"

	"github.com/dylEasydev/go-oauth2-easyclass/validators"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// structure des information sur un client oauth2
type InfoClient struct {
	ID                  uuid.UUID `gorm:"primarykey;type:uuid;default:uuid_generate_v4()"`
	NameOrganization    string    `gorm:"not null;unique" validate:"required,name"`
	TypeApplication     string    `gorm:"not null" validate:"required,appallowed"`
	AddressOrganization string    `gorm:"not null;unique" validate:"required,email"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Image Image `gorm:"polymorphic:Picture;"`
}

func (InfoClient) TableName() string {
	return "info_clients"
}

// validation du Model avant avant la sauvegarde
func (info *InfoClient) BeforeSave(tx *gorm.DB) error {
	return validators.ValidateStruct(info)
}
