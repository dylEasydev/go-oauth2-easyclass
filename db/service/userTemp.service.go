package service

import (
	"context"
	"errors"

	"github.com/dylEasydev/go-oauth2-easyclass/db/interfaces"
	"gorm.io/gorm"
)

var ErrDestroy = errors.New("impossible de supprimer l'utilisateur temporaire")

type UsertempService struct {
	Ctx *context.Context
	Db  *gorm.DB
}

func (service *UsertempService) SaveUser(user interfaces.UserTempInterafce) error {
	if err := user.SavePerm(service.Db); err != nil {
		return err
	}
	if err := user.DestroyUser(service.Db); err != nil {
		return ErrDestroy
	}
	return nil
}
