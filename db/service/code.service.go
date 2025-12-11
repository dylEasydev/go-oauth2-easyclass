package service

import (
	"context"
	"errors"

	"github.com/dylEasydev/go-oauth2-easyclass/db/interfaces"
	"github.com/dylEasydev/go-oauth2-easyclass/db/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrNotCode = errors.New("mauvais code de verification")

type CodeService struct {
	Ctx *context.Context
	Db  *gorm.DB
}

func (service *CodeService) FindCode(code string, id uuid.UUID) (*models.CodeVerif, error) {
	codeVerif, err := gorm.G[models.CodeVerif](service.Db).Where(&models.CodeVerif{Code: code, VerifiableID: id}).First(*service.Ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotCode
		}
		return nil, err
	}
	return &codeVerif, nil
}

func (service *CodeService) UpdateCodeVerif(user interfaces.UserInterface, code string) error {
	beforeCode, err := service.FindCode(code, user.GetId())
	if beforeCode == nil {
		return err
	}
	_, err = gorm.G[models.CodeVerif](service.Db).Where("id = ?", beforeCode.ID).Updates(*service.Ctx, *beforeCode)
	if err != nil {
		return err
	}
	return nil
}
