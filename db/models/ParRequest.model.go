package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type PARRequest struct {
	ID         uuid.UUID      `gorm:"primarykey;type:uuid;default:uuid_generate_v4()"`
	RequestURI string         `gorm:"unique;not null"`
	Form       datatypes.JSON `gorm:"type:jsonb"`
	ExpiresAt  time.Time      `gorm:"index"`
	Used       bool           `gorm:"default:false;"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	ClientID  uuid.UUID `gorm:"type:uuid;not null"`
	Client    Client    `gorm:"foreignKey:ClientID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	SessionID uuid.UUID `gorm:"type:uuid;not null"`
	Session   Session   `gorm:"foreignKey:SessionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (PARRequest) TableName() string {
	return "par_requests"
}
