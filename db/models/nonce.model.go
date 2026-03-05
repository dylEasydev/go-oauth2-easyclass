package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// structure Nonce OIDC
// pour l'extention **verfiable**
type Nonce struct {
	ID uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`

	//jeton d'accès
	AccessToken string    `gorm:"not null"`
	Nonce       string    `gorm:"unique;not null;index"`
	ExpiresAt   time.Time `gorm:"type:timestamptz;index"`

	//timestamps
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// implementation de l'interface Tabler
func (Nonce) TableName() string {
	return "nonces"
}
