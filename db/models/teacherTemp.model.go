package models

import (
	"fmt"

	"github.com/dylEasydev/go-oauth2-easyclass/validators"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// structure des ensignant temporaire
type TeacherTemp struct {
	TeacherBase
	//code de verification envoyer par mail pour sa validation
	CodeVerif CodeVerif `gorm:"polymorphic:Verifiable;"`
}

func (TeacherTemp) TableName() string {
	return "teacher_temp"
}

func (teacher *TeacherTemp) BeforeSave(tx *gorm.DB) (err error) {
	if err = validators.ValidateStruct(teacher); err != nil {
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(teacher.Password), 10)
	if err != nil {
		return
	}

	teacher.Password = string(hash)
	return
}

func (teacher *TeacherTemp) AfterCreate(tx *gorm.DB) (err error) {
	code := CodeVerif{
		VerifiableID:   teacher.ID,
		VerifiableType: teacher.TableName(),
	}

	if err = tx.Create(&code).Error; err != nil {
		return err
	}

	return nil
}

func (teacher *TeacherTemp) SavePerm(tx *gorm.DB) error {
	teacherWait := TeacherWaiting{
		TeacherBase: TeacherBase{
			UserBase: UserBase{
				UserName: teacher.UserName,
				Email:    teacher.Email,
				Password: teacher.Password,
			},
			SubjectName: teacher.SubjectName,
		},
	}
	if err := tx.Create(&teacherWait).Error; err != nil {
		return fmt.Errorf("erreur de création de l'enseignat en attente: %w", err)
	}
	return nil
}
