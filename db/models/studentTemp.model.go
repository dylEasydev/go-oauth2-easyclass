package models

import (
	"fmt"

	"github.com/dylEasydev/go-oauth2-easyclass/db/query"
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

// hooks avant la sauvegarde
// validation et hash du password
func (student *StudentTemp) BeforeSave(tx *gorm.DB) error {
	if err := validators.ValidateStruct(student); err != nil {
		return err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(student.Password), Cout_hash)
	if err != nil {
		return err
	}

	student.Password = string(hash)
	return nil
}

// hooks après la creation de l'utilisateur
// création du code de verfication et envoie par mail
func (student *StudentTemp) AfterCreate(tx *gorm.DB) (err error) {
	code := CodeVerif{
		VerifiableID:   student.ID,
		VerifiableType: student.TableName(),
	}

	if err = query.QueryCreate(tx, &code); err != nil {
		return err
	}

	return nil
}

// sauvegarde de l'utilisateur en tant qu'utilisteur permanent
func (student *StudentTemp) SavePerm(tx *gorm.DB) (err error) {
	//initialisation de la transaction
	return tx.Transaction(func(tx *gorm.DB) error {
		// session bd sans hooks
		txhooks := tx.Session(&gorm.Session{SkipHooks: true})

		// association de l'utilisateur à un role
		role := Role{
			RoleName:     "student",
			RoleDescript: "role de l'étudiant",
		}
		if err := tx.Where(Role{RoleName: role.RoleName}).FirstOrCreate(&role).Error; err != nil {
			return fmt.Errorf("erreur lors de la création du rôle: %w", err)
		}

		user := User{
			UserBase: UserBase{
				UserName: student.UserName,
				Email:    student.Email,
				Password: student.Password,
			},
			RoleID: role.ID,
			Role:   role,
			Image: Image{
				PicturesName: "profil_default.png",
				UrlPictures:  fmt.Sprintf("%s/public/profil_default.png", utils.URL_Image),
			},
		}

		// création de l'utilisateur permanent
		if err := query.QueryCreate(txhooks, &user); err != nil {
			return fmt.Errorf("erreur lors de la création de l'utilisateur: %w", err)
		}

		// création du code de vérification
		code := CodeVerif{
			VerifiableID:   user.ID,
			VerifiableType: user.TableName(),
		}
		err := code.BeforeSave(tx)
		if err != nil {
			return fmt.Errorf("erreur lors de la création du code de vérification: %w", err)
		}
		if err := query.QueryCreate(txhooks, &code); err != nil {
			return fmt.Errorf("erreur lors de la création du code de vérification: %w", err)
		}

		return nil
	})
}

func (student *StudentTemp) DestroyUser(tx *gorm.DB) error {
	return query.QueryDeleteById[StudentTemp](tx, student.ID)
}
