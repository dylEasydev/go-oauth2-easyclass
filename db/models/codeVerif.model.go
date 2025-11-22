package models

import (
	"context"
	"fmt"
	"log"
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

// structure du code de verificatioon
// util pour valider le mail fournir par l'utilisateur
type CodeVerif struct {
	ID uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	//code hashé en BD
	Code string `gorm:"not null"`
	//code en clair
	rawCode string

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
}

// implementation de l'interface Tabler
func (CodeVerif) TableName() string {
	return "code_verifs"
}

// récupération de l'objet polymorphe scanner dans la base de leur héritage
func (codeverif *CodeVerif) GetForeign(tx *gorm.DB) (interfaces.UserInterface, error) {
	foreign := UserBase{}
	if err := tx.Table(codeverif.VerifiableType).Select("id", "email", "username").Where(map[string]any{"id": codeverif.VerifiableID}).Take(&foreign).Error; err != nil {
		return nil, err
	}
	return &foreign, nil
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
		codeverif.ExpiresAt = time.Now().Add(CODE_VALIDATE).UTC()
	}
	return nil
}

// envoye du mail après créer
func (codeVerif *CodeVerif) AfterSave(tx *gorm.DB) (err error) {

	//récupération de l'objet polymorphes
	verifiable, err := codeVerif.GetForeign(tx)
	if err != nil {
		return
	}

	email := verifiable.GetMail()
	code := codeVerif.rawCode
	name := verifiable.GetName()

	go func(email string, name string, code string) {
		if mailErr := utils.SendVerificationCode(email, code, name); mailErr != nil {
			log.Printf("warning: failed to send verification email to %s: %v", email, mailErr)
		} else {
			log.Printf("info: verification email sent to %s", email)
		}
	}(email, name, code)

	return nil
}

// verification de l'expiration
func (codeverif *CodeVerif) IsExpired() bool {
	return time.Now().UTC().Before(codeverif.ExpiresAt)
}

// verifier s'il est déjà utiliser
func (codeVerif *CodeVerif) IsUsed() bool {
	return codeVerif.UseAt != nil
}

// marké un code déjà utiliser
func (codeverif *CodeVerif) MarkUsed(tx *gorm.DB) error {
	txSession := tx.Session(&gorm.Session{SkipHooks: true})
	ctx := context.Background()
	now := time.Now().UTC()
	_, err := gorm.G[CodeVerif](txSession).Where(map[string]any{"id": codeverif.ID}).Updates(ctx, CodeVerif{UseAt: &now})
	return err
}
