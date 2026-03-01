package db

//pacakeges db

import (
	"errors"
	"fmt"
	"log"
	"os"

	jose "github.com/go-jose/go-jose/v3"

	"github.com/dylEasydev/go-oauth2-easyclass/db/models"
	"github.com/dylEasydev/go-oauth2-easyclass/db/query"
	"github.com/dylEasydev/go-oauth2-easyclass/utils"
	"github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// structure de sauvegarde
type Store struct {
	db *gorm.DB
}

func (store *Store) GetDb() *gorm.DB {
	return store.db
}

// initialisation de la DB
func InitDB(db *gorm.DB) error {
	//session de BD avec hooks
	txhooks := db.Session(&gorm.Session{SkipHooks: true})

	//instance de db avec clause on Donothing à true
	txClause := db.Clauses(clause.OnConflict{DoNothing: true})

	//lecture de l'ensemble des permisssions
	name := "scope_app"
	data, err := utils.ReadJSON[models.ScopeData](name)
	if err != nil {
		return err
	}

	//création de l'emsble des scopes pour le bulk create
	scopes := make([]models.Scope, 0, len(data.Data))
	for _, elem := range data.Data {
		scopes = append(scopes, models.Scope{
			ScopeName:     elem.ScopeName,
			ScopeDescript: elem.ScopeDescript,
		})
	}

	if err := txClause.Create(&scopes).Error; err != nil {
		return fmt.Errorf("erreur lors de la création de scopes: %w", err)
	}

	//création du roles de l'administrateur
	role := models.Role{RoleName: "admin", RoleDescript: "role de l'administrateur"}
	if err = db.Where(models.Role{RoleName: role.RoleName}).FirstOrCreate(&role).Error; err != nil {
		return err
	}

	//creation de l'utilisateur administrateur
	username := os.Getenv("USER_NAME")
	email := os.Getenv("COMPANING_MAIl")
	password := os.Getenv("USER_PASSWORD")

	//recherche si l'utlisateur existe dejà si non le créer
	var user models.User
	if err := db.Where("user_name = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			//si l'utilsateur n'existe pas
			user = models.User{
				UserBase: models.UserBase{
					UserName: username,
					Email:    email,
					Password: password,
				},
				Role:   role,
				RoleID: role.ID,
				Image: models.Image{
					PicturesName: "profil_default.png",
					UrlPictures:  fmt.Sprintf("%s/public/profil_default.png", utils.URL_Image),
				},
			}
			if err := query.QueryCreate(db, &user); err != nil {
				return err
			}

			//création de son code de verification
			code := models.CodeVerif{
				VerifiableID:   user.ID,
				VerifiableType: user.TableName(),
			}

			//validation et  hashage du code de verfication
			if err := code.BeforeSave(db); err != nil {
				return fmt.Errorf("erreur lors de la création du code de vérification: %w", err)
			}
			if err := query.QueryCreate(txhooks, &code); err != nil {
				return fmt.Errorf("erreur lors de la création du code de vérification: %w", err)
			}
		} else {
			return err
		}
	}

	//création des informations du client
	compEmail := os.Getenv("COMPANING_MAIl")
	info := models.InfoClient{
		NameOrganization:    "EasyClassOrg",
		TypeApplication:     "web app",
		AddressOrganization: compEmail,
		Image: models.Image{
			PicturesName: "client_default.png",
			UrlPictures:  fmt.Sprintf("%s/public/client_default.png", utils.URL_Image),
		},
	}
	if err = db.Where(models.InfoClient{NameOrganization: info.NameOrganization}).FirstOrCreate(&info).Error; err != nil {
		return err
	}

	//création du client oidc
	client := models.Client{
		Active:                  utils.PtrBool(true),
		Secret:                  os.Getenv("SECRET_CLIENT"),
		RotatedSecrets:          pq.StringArray{os.Getenv("SECRET_CLIENT2")},
		Public:                  utils.PtrBool(false),
		RedirectURIs:            pq.StringArray{"https://localhost:3000/callback", "https://127.0.0.1:3000/callback"},
		Scopes:                  pq.StringArray{"openid", "admin.*"},
		Audience:                pq.StringArray{},
		Grants:                  pq.StringArray{"code", "token", "client_credentials", "password"},
		ResponseTypes:           pq.StringArray{"code", "token"},
		RequestURIs:             pq.StringArray{},
		ResponseModes:           pq.StringArray{"query", "fragment", "form_post"},
		TokenEndpointAuthMethod: "client_secret_basic",
		InfoClientID:            info.ID,
		InfoClient:              info,
	}

	if err = db.Where(models.Client{InfoClientID: info.ID}).FirstOrCreate(&client).Error; err != nil {
		return err
	}

	//création de la clé de validations du client
	var existingKey models.ClientKey
	if err = db.Where("client_id = ?", client.ID).First(&existingKey).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			//chargement de la clé RSA du client
			secret, err := utils.LoadPublicKey("public")
			if err != nil {
				return err
			}

			jwk := jose.JSONWebKey{
				Key:       secret,
				KeyID:     "init-key",
				Algorithm: "RS256",
				Use:       "sig",
			}
			ck := models.ClientKey{
				Issuer:    client.ID.String(),
				Subject:   client.ID.String(),
				KeyID:     "init-key",
				Algorithm: "RS256",
				Scopes:    client.Scopes,
				JWK:       models.JWKey(jwk),
				ClientID:  client.ID,
			}
			if err := query.QueryCreate(db, &ck); err != nil {
				return err
			}
		} else {
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
		//Logger:      logger.Default.LogMode(logger.Info),
		PrepareStmt: true,
	})
	if err != nil {
		log.Fatalf("erreur de connexion à la base de données: %v", err)
	}

	// s'assurer que  uuid_generate_v4() exist (uuid-ossp extension) avant de lancer les migrations
	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`).Error; err != nil {
		log.Fatalf("erreur de création de l'extension uuid-ossp: %v", err)
	}

	//migration de la BD
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
		models.AuthPermission{},
	)

	if err != nil {
		log.Fatal("erreur de démarrage de la migrations:", err)
	}

	if err := InitDB(db); err != nil {
		log.Fatal("initialisation de la BD failed:", err)
	}
	return &Store{
		db: db,
	}
}
