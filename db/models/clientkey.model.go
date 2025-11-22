package models

//packages models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

//db models ClienKey pour OIDC
//client lié à la norme rfc7523
//le client fournis un jeton assertion
// il peremet de stocké la clé public de la clé privé
//que le client utilise pour générer le jwt assertion

type ClientKey struct {
	ID        uuid.UUID      `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Issuer    string         `gorm:"not null;uniqueIndex:idx_issuer_subject"`
	Subject   string         `gorm:"not null;uniqueIndex:idx_issuer_subject"`
	KeyID     string         `gorm:"not null"`
	Algorithm string         `gorm:"not null"`
	Scopes    pq.StringArray `gorm:"type:text[]"`
	JWK       JWKey          `gorm:"type:jsonb;not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	//raltion avec le client assoccier à la clé
	ClientID uuid.UUID `gorm:"type:uuid;not null;"`
	Client   Client    `gorm:"foreignKey:ClientID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// implementation de interface Tabler(pour le nom de la table)
func (ClientKey) TableName() string {
	return "client_keys"
}
