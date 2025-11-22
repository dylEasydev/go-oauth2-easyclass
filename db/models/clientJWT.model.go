package models

//packages models

import (
	"time"

	"github.com/dylEasydev/go-oauth2-easyclass/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// db models clientJWT pour OIDC
type ClientJWT struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Active    *bool     `gorm:"default:true"`
	JTI       string    `gorm:"unique;not null"`
	ExpiresAt time.Time

	//timestamps
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// implementation de interface Tabler(pour le nom de la table)
func (ClientJWT) TableName() string {
	return "client_jwts"
}

// hooks avant la sauvegarde
func (c *ClientJWT) BeforeCreate(tx *gorm.DB) (err error) {
	if c.Active == nil {
		c.Active = utils.PtrBool(true)
	}
	return nil
}

// verifie si le clientJWt est expiré
func (c *ClientJWT) IsExpired() bool {
	return time.Now().UTC().Before(c.ExpiresAt)
}

// verifie si le clientJWt estr encore valide
func (c *ClientJWT) IsValid() bool {
	return c.Active != nil && *c.Active && !c.IsExpired()
}
