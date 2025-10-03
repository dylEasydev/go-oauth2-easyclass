package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// strcture du jeton de rafraichissement d'authentification oauth2
type RefreshToken struct {
	ID     uuid.UUID `gorm:"primarykey;type:uuid;default:uuid_generate_v4()"`
	Active *bool     `gorm:"default:true"`

	//valeur du token de rafraichissement
	Signature string `gorm:"unique;not null"`

	RequestedAt time.Time `gorm:"type:jsonb;default:null"`

	//permission et grant_types demandés dans la requêtes
	RequestedScopes pq.StringArray `gorm:"type:text[]"`
	GrantedScopes   pq.StringArray `gorm:"type:text[]"`

	Form datatypes.JSON

	//permissions et grant_types accordé
	RequestedAudience pq.StringArray `gorm:"type:text[]"`
	GrantedAudience   pq.StringArray `gorm:"type:text[]"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	//realtion entre refreshtoken , session et client
	ClientID  uuid.UUID  `gorm:"type:uuid;not null"`
	Client    Client     `gorm:"foreignKey:ClientID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	SessionID *uuid.UUID `gorm:"type:uuid"`
	Session   Session    `gorm:"foreignKey:SessionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}
