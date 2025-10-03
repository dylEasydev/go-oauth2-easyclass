package models

import (
	"time"

	"github.com/dylEasydev/go-oauth2-easyclass/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ClientJWT struct {
	ID        uuid.UUID `gorm:"primarykey;type:uuid;default:uuid_generate_v4()"`
	Active    *bool     `gorm:"default:true"`
	JTI       string    `gorm:"unique;not null"`
	ExpiresAt time.Time

	//timestamps
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	ClientID uuid.UUID `gorm:"type:uuid;not null"`
	Client   Client    `gorm:"foreignKey:ClientID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (ClientJWT) TableName() string {
	return "client_jwts"
}

func (c *ClientJWT) BeforeCreate(tx *gorm.DB) (err error) {
	if c.Active == nil {
		c.Active = utils.PtrBool(true)
	}
	return nil
}

func (c *ClientJWT) IsExpired() bool {
	return time.Now().After(c.ExpiresAt)
}

func (c *ClientJWT) IsValid() bool {
	return c.Active != nil && *c.Active && !c.IsExpired()
}
