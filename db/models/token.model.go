package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// strcture du Jeton d'access d'authentification oauth2
type AccessToken struct {
	ID     uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Active *bool     `gorm:"default:true"`

	//valeur unique du token ( bloc de carractère générer aléatoirement)
	Signature string `gorm:"uniqueIndex;not null"`

	RequestedAt time.Time

	//permission et grant_types demandés dans la requêtes
	RequestedScopes pq.StringArray `gorm:"type:text[]"`
	GrantedScopes   pq.StringArray `gorm:"type:text[]"`

	Form datatypes.JSON `gorm:"type:jsonb;default:null"`

	//permissions et grant_types accordé
	RequestedAudience pq.StringArray `gorm:"type:text[]"`
	GrantedAudience   pq.StringArray `gorm:"type:text[]"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// relation hasone avec les models Client et session
	ClientID  uuid.UUID  `gorm:"type:uuid;not null"`
	Client    Client     `gorm:"foreignKey:ClientID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	SessionID *uuid.UUID `gorm:"type:uuid"`
	Session   Session    `gorm:"foreignKey:SessionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (AccessToken) TableName() string {
	return "access_tokens"
}
