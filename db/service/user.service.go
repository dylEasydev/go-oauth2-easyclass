package service

import (
	"context"
	"fmt"

	"github.com/dylEasydev/go-oauth2-easyclass/db/models"
	"github.com/dylEasydev/go-oauth2-easyclass/db/query"
	"github.com/dylEasydev/go-oauth2-easyclass/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserService struct {
	Ctx context.Context
	Db  *gorm.DB
}
type UserBody struct {
	Name     string
	Email    string
	Password string
}

func InitUserService(ctx context.Context, db *gorm.DB) *UserService {
	return &UserService{
		Ctx: ctx,
		Db:  db,
	}
}

func (service *UserService) FindUserByName(name, email string) (*models.User, error) {

	user, err := gorm.G[models.User](service.Db).Joins(clause.JoinTarget{Association: "CodeVerif"}, nil).Joins(clause.JoinTarget{Association: "Image"}, nil).Joins(clause.JoinTarget{Association: "Role"}, nil).Preload("Role.Scopes", nil).Where("user_name = ? or email = ?", name, email).First(service.Ctx)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (service *UserService) FindUserById(id uuid.UUID) (*models.User, error) {
	user, err := gorm.G[models.User](service.Db).Joins(clause.JoinTarget{Association: "CodeVerif"}, nil).Joins(clause.JoinTarget{Association: "Image"}, nil).Joins(clause.JoinTarget{Association: "Role"}, nil).Preload("Role.Scopes", nil).Where("id = ? ", id).First(service.Ctx)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (service *UserService) CreateUser(data *UserBody) (*models.User, error) {
	var newUser *models.User
	err := service.Db.WithContext(service.Ctx).Transaction(func(tx *gorm.DB) error {
		txhooks := tx.Session(&gorm.Session{SkipHooks: true})
		// association de l'utilisateur à un role
		role := models.Role{
			RoleName:     "admin",
			RoleDescript: "role de l'administrateur",
		}
		if err := tx.Where(models.Role{RoleName: role.RoleName}).FirstOrCreate(&role).Error; err != nil {
			return fmt.Errorf("erreur lors de la création du rôle: %w", err)
		}

		//création de l'utilisateur
		newUser = &models.User{
			UserBase: models.UserBase{
				UserName: data.Name,
				Email:    data.Email,
				Password: data.Password,
			},
			RoleID: role.ID,
			Role:   role,
			Image: models.Image{
				PicturesName: "profil_default.png",
				UrlPictures:  fmt.Sprintf("%s/public/profil_default.png", utils.URL_Image),
			},
		}

		if err := query.QueryCreate(tx, newUser); err != nil {
			return fmt.Errorf("erreur lors de la création de l'utilisateur: %w", err)
		}

		// création du code de vérification
		code := models.CodeVerif{
			VerifiableID:   newUser.ID,
			VerifiableType: newUser.TableName(),
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

	return newUser, err
}
