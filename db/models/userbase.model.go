package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// structure de model d'utilisateur de base
type UserBase struct {
	ID       uuid.UUID `gorm:"primaryey;type:uuid;default:uuid_generate_v4()"`
	UserName string    `gorm:"column:username;not null;unique" validate:"required,name"`
	Password string    `gorm:"column:password;not null" validate:"required,min=8,password"`
	Email    string    `gorm:"column:email;unique" validate:"required,email"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// implementation de interfac userInterface
func (user *UserBase) GetMail() string {
	return user.Email
}

func (user *UserBase) GetName() string {
	return user.UserName
}
