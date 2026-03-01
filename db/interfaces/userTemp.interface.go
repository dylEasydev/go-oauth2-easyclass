package interfaces

import "gorm.io/gorm"

type UserTempInterface interface {
	UserInterface
	SavePerm(tx *gorm.DB) error
	DestroyUser(tx *gorm.DB) error
}
