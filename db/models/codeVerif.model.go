package models

import (
	"fmt"
	"time"

	"github.com/dylEasydev/go-oauth2-easyclass/db/interfaces"
	"github.com/dylEasydev/go-oauth2-easyclass/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	CODE_LENGTH   = 6
	CODE_VALIDATE = 1 * time.Hour
)

// structure du code de verification
type CodeVerif struct {
	ID uuid.UUID `gorm:"primarykey;type:uuid;default:uuid_generate_v4()"`
	//code hashé en BD
	Code string
	//code en clair
	rawCode string `gorm:"-"`

	//temps d'expiration du code de verification
	ExpiresAt time.Time

	//temps d'utilisation
	UseAt *time.Time

	//timestamps
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	VerifiableID   uuid.UUID `gorm:"not null;"`
	VerifiableType string    `gorm:"not null;"`

	Verifiable interfaces.UserInterface `gorm:"polymorphic:Verifiable;"`
}

// implementation de l'interface Tabler
func (CodeVerif) TableName() string {
	return "code_verifs"
}

// validation du model avant la sauvegarde
func (codeverif *CodeVerif) BeforeSave(tx *gorm.DB) (err error) {
	raw, err := utils.GenerateVerificationCode()
	if err != nil {
		return fmt.Errorf("erreur génération code: %w", err)
	}
	codeverif.rawCode = raw
	codeverif.Code = utils.GenerateHash(raw)

	if codeverif.ExpiresAt.IsZero() {
		codeverif.ExpiresAt = time.Now().Add(CODE_VALIDATE)
	}
	return nil
}

// envoye du mail après créer
func (codeVerif *CodeVerif) AfterSave(tx *gorm.DB) (err error) {
	if err = tx.Model(codeVerif).Association("Verifiable").Find(&codeVerif.Verifiable); err != nil {
		return
	}
	email := codeVerif.Verifiable.GetMail()
	code := codeVerif.rawCode

	go func(email string, code string) {
		mailErr := utils.SendVerificationCode(email, code)
		if mailErr != nil {
			fmt.Printf("problème lors de l'envoie de mail :%s\n", mailErr)
		} else {
			fmt.Println("mail bien envoyé")
		}
	}(email, code)

	return nil
}

// verification de l'expiration
func (codeverif *CodeVerif) IsExpired() bool {
	return time.Now().After(codeverif.ExpiresAt)
}

// verifier s'il est déjà utiliser
func (codeVerif *CodeVerif) IsUsed() bool {
	return codeVerif.UseAt != nil
}

// marké un code déjà utiliser
func (codeverif *CodeVerif) MarkUsed(tx *gorm.DB) error {
	now := time.Now()
	codeverif.UseAt = &now
	return tx.Session(&gorm.Session{SkipHooks: true}).Save(codeverif).Error
}
