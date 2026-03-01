package service

import (
	"github.com/dylEasydev/go-oauth2-easyclass/db/interfaces"
	"gorm.io/gorm"
)

type UsertempService struct {
	Db *gorm.DB
}

func InitUserTempService(db *gorm.DB) *UsertempService {
	return &UsertempService{
		Db: db,
	}
}

// sauvegarde de user_temp en user
func (service *UsertempService) SaveUser(user interfaces.UserTempInterface) error {
	if err := user.SavePerm(service.Db); err != nil {
		return err
	}
	if err := user.DestroyUser(service.Db); err != nil {
		return ErrDestroy
	}
	return nil
}
