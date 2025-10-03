package models

import (
	"github.com/dylEasydev/go-oauth2-easyclass/validators"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
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

func (user *User) BeforeSave(tx *gorm.DB) (err error) {
	if err = validators.ValidateStruct(user); err != nil {
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return
	}

	user.Password = string(hash)
	return
}
