package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

//models PAR OIDC
// pour l'extension PAR request

type PARRequest struct {
	ID         uuid.UUID      `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	RequestURI string         `gorm:"uniqueIndex;not null"`
	Form       datatypes.JSON `gorm:"type:jsonb"`
	ExpiresAt  time.Time      `gorm:"index"`
	Used       bool           `gorm:"default:false;"`

	//Permissions et Grant demandés
	RequestedScopes pq.StringArray `gorm:"type:text[]"`
	GrantedScopes   pq.StringArray `gorm:"type:text[]"`

	//Permissions et grant acceptés
	RequestedAudience pq.StringArray `gorm:"type:text[]"`
	GrantedAudience   pq.StringArray `gorm:"type:text[]"`

	RequestedAt time.Time

	RedirectURI  datatypes.JSON `gorm:"type:jsonb"`
	ResponseMode string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	ClientID  uuid.UUID `gorm:"type:uuid;not null"`
	Client    Client    `gorm:"foreignKey:ClientID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	SessionID uuid.UUID `gorm:"type:uuid;not null"`
	Session   Session   `gorm:"foreignKey:SessionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// implementation de l'interface Tabler
func (PARRequest) TableName() string {
	return "par_requests"
}
