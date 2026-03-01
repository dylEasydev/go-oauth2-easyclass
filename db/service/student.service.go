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

func InitStudentService(ctx *context.Context, db *gorm.DB) *StudentService {
	return &StudentService{
		Ctx: ctx,
		Db:  db,
	}
}

func (service *StudentService) CreateUser(data *UserBody) (*models.StudentTemp, error) {
	studentFind, err := FindUserByName[models.StudentTemp](service.Ctx, service.Db, data.Name, data.Email)
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
			return &newStudent, nil
		}
		return nil, err
	}
	studentUpdate := models.StudentTemp{
		UserBase: models.UserBase{
			ID:       studentFind.ID,
			UserName: data.Name,
			Email:    data.Email,
			Password: data.Password,
		},
	}
	err = service.Db.WithContext(*service.Ctx).Model(studentFind).Where("id = ?", studentFind.ID).Updates(&studentUpdate).Error
	if err != nil {
		return nil, err
	}
	codeservice := InitCodeService(service.Ctx, service.Db)
	err = codeservice.UpdateCodeVerif(studentFind, studentFind.CodeVerif.Code)
	if err != nil {
		return nil, err
	}
	return &studentUpdate, nil
}
