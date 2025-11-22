package models

// packages models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// db models d'AuthorizationCode pour OIDC
type AuthorizationCode struct {
	ID     uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Active *bool     `gorm:"default:true"`

	//code d'authorization unique
	Code string `gorm:"unique;not null"`

	//temps d'émission de la requête
	RequestedAt time.Time

	//Permissions et Grant demandés
	RequestedScopes pq.StringArray `gorm:"type:text[]"`
	GrantedScopes   pq.StringArray `gorm:"type:text[]"`

	//formulaire d'authorize request
	Form datatypes.JSON `gorm:"type:jsonb;default:null"`

	//Permissions et grant acceptés
	RequestedAudience pq.StringArray `gorm:"type:text[]"`
	GrantedAudience   pq.StringArray `gorm:"type:text[]"`

	//timestamps
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	//realtion avec client et session
	ClientID  uuid.UUID `gorm:"type:uuid;not null"`
	Client    Client    `gorm:"foreignKey:ClientID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	SessionID uuid.UUID `gorm:"type:uuid;not null"`
	Session   Session   `gorm:"foreignKey:SessionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// implementation de interface Tabler(pour le nom de la table)
func (AuthorizationCode) TableName() string {
	return "authorization_codes"
}
