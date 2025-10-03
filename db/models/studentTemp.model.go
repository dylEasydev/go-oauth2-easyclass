package models

import (
	"fmt"

	"github.com/dylEasydev/go-oauth2-easyclass/utils"
	"github.com/dylEasydev/go-oauth2-easyclass/validators"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// structure dest étudiants temporaires
type StudentTemp struct {
	UserBase

	//code de verification envoyer par mail pour sa validation
	CodeVerif CodeVerif `gorm:"polymorphic:Verifiable;"`
}

func (StudentTemp) TableName() string {
	return "student_temps"
}
func (student *StudentTemp) BeforeSave(tx *gorm.DB) (err error) {
	if err = validators.ValidateStruct(student); err != nil {
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(student.Password), 10)
	if err != nil {
		return
	}

	student.Password = string(hash)
	return
}

func (student *StudentTemp) AfterCreate(tx *gorm.DB) (err error) {
	code := CodeVerif{
		VerifiableID:   student.ID,
		VerifiableType: student.TableName(),
	}

	if err = tx.Create(&code).Error; err != nil {
		return err
	}

	return nil
}

func (student *StudentTemp) SavePerm(tx *gorm.DB) (err error) {
	return tx.Transaction(func(tx *gorm.DB) error {
		txhooks := tx.Session(&gorm.Session{SkipHooks: true})

		role := Role{
			RoleDescript: "role de l'étudiant",
		}
		if err := tx.FirstOrCreate(&role, Role{RoleName: "student"}).Error; err != nil {
			return fmt.Errorf("erreur lors de la création du rôle: %w", err)
		}

		user := User{
			UserBase: UserBase{
				UserName: student.UserName,
				Email:    student.Email,
				Password: student.Password,
			},
			RoleID: role.ID,
		}

		// création de l'utilisateur permanent
		if err := txhooks.Create(&user).Error; err != nil {
			return fmt.Errorf("erreur lors de la création de l'utilisateur: %w", err)
		}

		image := Image{
			PicturesName: "profil_default.png",
			UrlPictures:  fmt.Sprintf("%s/public/profil_default.png", utils.URL_Image),
			PictureID:    user.ID,
			PictureType:  user.TableName(),
		}

		// création de l'image de profil
		if err := tx.Create(&image).Error; err != nil {
			return fmt.Errorf("erreur lors de la création de l'image de profil: %w", err)
		}

		// création du code de vérification
		code := CodeVerif{
			VerifiableID:   user.ID,
			VerifiableType: user.TableName(),
		}

		if err := txhooks.Create(&code).Error; err != nil {
			return fmt.Errorf("erreur lors de la création du code de vérification: %w", err)
		}

		return nil
	})
}
