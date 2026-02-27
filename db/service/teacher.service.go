package service

import (
	"context"
	"errors"

	"github.com/dylEasydev/go-oauth2-easyclass/db/models"
	"github.com/dylEasydev/go-oauth2-easyclass/db/query"
	"gorm.io/gorm"
)

type TeacherService struct {
	Ctx context.Context
	Db  *gorm.DB
}

type TeacherBody struct {
	UserBody
	Subject string
}

func InitTeacherService(ctx context.Context, db *gorm.DB) *TeacherService {
	return &TeacherService{
		Ctx: ctx,
		Db:  db,
	}
}

func (service *TeacherService) CreateUser(data *TeacherBody) (*models.TeacherTemp, error) {
	teacherFind, err := FindUserByName[models.TeacherTemp](service.Ctx, service.Db, data.Name, data.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newTeacher := models.TeacherTemp{
				TeacherBase: models.TeacherBase{
					UserBase: models.UserBase{
						UserName: data.Name,
						Email:    data.Email,
						Password: data.Password,
					},
					SubjectName: data.Subject,
				},
			}
			if err := query.QueryCreate(service.Db, &newTeacher); err != nil {
				return nil, err
			}
			return &newTeacher, err
		}
		return nil, err
	}
	teacherUpadate := models.TeacherTemp{
		TeacherBase: models.TeacherBase{
			UserBase: models.UserBase{
				UserName: data.Name,
				Email:    data.Email,
				Password: data.Password,
			},
			SubjectName: data.Subject,
		},
	}
	_, err = gorm.G[models.TeacherTemp](service.Db).Where("id = ?", teacherFind.ID).Updates(service.Ctx, teacherUpadate)
	if err != nil {
		return nil, err
	}
	codeservice := InitCodeService(service.Ctx, service.Db)
	err = codeservice.UpdateCodeVerif(teacherFind, teacherFind.CodeVerif.Code)
	if err != nil {
		return nil, err
	}
	return &teacherUpadate, nil
}
