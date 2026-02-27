package models

import (
	"github.com/dylEasydev/go-oauth2-easyclass/validators"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// constante de hash bcrypt
const (
	Cout_hash = 10
)

// structure du model utilisateur permanent
type User struct {
	UserBase

	CodeVerif CodeVerif `gorm:"polymorphic:Verifiable;"`
	Image     Image     `gorm:"polymorphic:Picture;"`

	//rôle de l'utilisateur (admin , student , teacher ...)
	Role   Role      `gorm:"foreignKey:RoleID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	RoleID uuid.UUID `gorm:"type:uuid;not null"`
}

func (User) TableName() string {
	return "user"
}

// si modification du password hash du mots de passe avant la sauvegarde
func (user *User) BeforeSave(tx *gorm.DB) error {
	if err := validators.ValidateStruct(user); err != nil {
		return err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), Cout_hash)
	if err != nil {
		return err
	}

	user.Password = string(hash)
	return nil
}
