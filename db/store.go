package db

//pacakeges db

import (
	"fmt"
	"log"
	"os"

	"github.com/dylEasydev/go-oauth2-easyclass/db/models"
	"github.com/dylEasydev/go-oauth2-easyclass/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// structure de sauvegarde
type Store struct {
	db *gorm.DB
}

// initialisation de la DB
func InitDB(db *gorm.DB) error {
	//création des permissions
	name := "ressources/scope_app"
	data, err := utils.ReadJSON[models.ScopeData](name)
	if err != nil {
		return err
	}
	for _, elem := range data.Data {
		if err = db.Model(&models.Scope{}).FirstOrCreate(&elem, models.Scope{ScopeName: elem.ScopeName}).Error; err != nil {
			return err
		}
	}

	return nil
}

func New() *Store {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:      logger.Default.LogMode(logger.Error),
		PrepareStmt: true,
	})
	if err != nil {
		log.Fatalf("erreur de connexion à la base de données: %v", err)
	}

	err = db.AutoMigrate(
		models.User{},
		models.CodeVerif{},
		models.Image{},
		models.InfoClient{},
		models.Client{},
		models.ClientJWT{},
		models.AuthorizationCode{},
		models.RefreshToken{},
		models.AccessToken{},
		models.Role{},
		models.Scope{},
		models.Session{},
		models.PKCE{},
		models.ClientKey{},
		models.PARRequest{},
		models.Nonce{},
		models.StudentTemp{},
		models.TeacherTemp{},
		models.TeacherWaiting{},
	)

	if err != nil {
		log.Fatal("failed to run migrations:", err)
	}

	return &Store{
		db: db,
	}
}
