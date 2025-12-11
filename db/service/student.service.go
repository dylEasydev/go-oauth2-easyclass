package service

import (
	"context"
	"errors"

	"github.com/dylEasydev/go-oauth2-easyclass/db/models"
	"github.com/dylEasydev/go-oauth2-easyclass/db/query"
	"gorm.io/gorm"
)

type StudentService struct {
	Ctx *context.Context
	Db  *gorm.DB
}

func (service *StudentService) CreateUser(data *UserBody) (*models.StudentTemp, error) {
	studentFind, err := FindUserByName[models.StudentTemp](*service.Ctx, service.Db, data.Name, data.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newStudent := models.StudentTemp{
				UserBase: models.UserBase{
					UserName: data.Name,
					Email:    data.Email,
					Password: data.Password,
				},
			}
			if err := query.QueryCreate(service.Db, &newStudent); err != nil {
				return nil, err
			}
			return &newStudent, err
		}
		return nil, err
	}
	studentUpdate := models.StudentTemp{
		UserBase: models.UserBase{
			UserName: data.Name,
			Email:    data.Email,
			Password: data.Password,
		},
	}
	_, err = gorm.G[models.StudentTemp](service.Db).Where("id = ?", studentFind.ID).Updates(*service.Ctx, studentUpdate)
	if err != nil {
		return nil, err
	}
	codeservice := &CodeService{Db: service.Db, Ctx: service.Ctx}
	err = codeservice.UpdateCodeVerif(studentFind, studentFind.CodeVerif.Code)
	if err != nil {
		return nil, err
	}
	return &studentUpdate, nil
}
