package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Nonce struct {
	ID uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`

	AccessToken string    `gorm:"not null;"`
	Nonce       string    `gorm:"unique;not null;"`
	ExpiresAt   time.Time `gorm:"index"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (Nonce) TableName() string {
	return "nonces"
}
