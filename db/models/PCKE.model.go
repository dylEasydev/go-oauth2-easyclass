package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type PKCE struct {
	ID        uuid.UUID `gorm:"primarykey;type:uuid;default:uuid_generate_v4()"`
	Active    *bool     `gorm:"default:true"`
	Signature string    `gorm:"unique"`

	RequestedAt time.Time

	RequestedScopes pq.StringArray `gorm:"type:text[]"`
	GrantedScopes   pq.StringArray `gorm:"type:text[]"`

	Form datatypes.JSON `gorm:"type:jsonb;default:null"`

	RequestedAudience pq.StringArray `gorm:"type:text[]"`
	GrantedAudience   pq.StringArray `gorm:"type:text[]"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	ClientID  uuid.UUID  `gorm:"type:uuid;not null"`
	Client    Client     `gorm:"foreignKey:ClientID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	SessionID *uuid.UUID `gorm:"type:uuid"`
	Session   Session    `gorm:"foreignKey:SessionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// implementation de l'interface Tabler
func (PKCE) TableName() string {
	return "pkces"
}
