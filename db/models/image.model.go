package models

import (
	"time"

	"github.com/dylEasydev/go-oauth2-easyclass/validators"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// structure des images
type Image struct {
	ID           uuid.UUID `gorm:"primarykey;type:uuid;default:uuid_generate_v4()"`
	PicturesName string    `gorm:"not null;unique;default:profil_default.png" validate:"required,name"`
	UrlPictures  string    `gorm:"not null;unique" validate:"required,url"`

	PictureID   uuid.UUID `gorm:"not null;"`
	PictureType string    `gorm:"not null;"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// implementation de l'interface Tabler
func (Image) TableName() string {
	return "images"
}

// validation du Model avant avant la sauvegarde
func (image *Image) BeforeSave(tx *gorm.DB) error {
	return validators.ValidateStruct(image)
}
