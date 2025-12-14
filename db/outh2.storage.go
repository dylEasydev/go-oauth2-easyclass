package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	"github.com/dylEasydev/go-oauth2-easyclass/db/models"
	"github.com/dylEasydev/go-oauth2-easyclass/utils"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/ory/fosite"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

//implementation de l'interface de CoreStorage

// AuthorizationStorage
func (store *Store) CreateAuthorizeCodeSession(ctx context.Context, code string, request fosite.Requester) (err error) {
	parsedID, err := uuid.Parse(request.GetID())
	if err != nil {
		return fmt.Errorf("request id invalide: %w", err)
	}
	client := request.GetClient()

	clientID, err := uuid.Parse(client.GetID())
	if err != nil {
		return fmt.Errorf("client id invalide: %w", err)
	}

	//marshalling du formulaire
	form, err := json.Marshal(request.GetRequestForm())
	if err != nil {
		return fmt.Errorf("erreur de marshalling du authorize form : %w", err)
	}

	// conversion de la session fosite
	session := request.GetSession().(*models.Session)

	//eneregistremment de la session en BD
	if err = gorm.G[models.Session](store.db).Create(ctx, session); err != nil {
		return fmt.Errorf("erreur de création de la sesion pour authorize code: %w", err)
	}

	//initialisation du code d'authorization
	data := models.AuthorizationCode{
		ID:                parsedID,
		Active:            utils.PtrBool(true),
		Code:              code,
		RequestedAt:       request.GetRequestedAt().UTC(),
		ClientID:          clientID,
		RequestedScopes:   pq.StringArray(request.GetRequestedScopes()),
		GrantedScopes:     pq.StringArray(request.GetGrantedScopes()),
		Form:              form,
		SessionID:         session.ID,
		RequestedAudience: pq.StringArray(request.GetRequestedAudience()),
		GrantedAudience:   pq.StringArray(request.GetGrantedAudience()),
	}

	// enregistrement du code d'authorization en BD
	if err = gorm.G[models.AuthorizationCode](store.db).Create(ctx, &data); err != nil {
		return fmt.Errorf("erreur de création d'authorize code: %w", err)
	}

	return nil
}

func (store *Store) GetAuthorizeCodeSession(ctx context.Context, code string, session fosite.Session) (request fosite.Requester, err error) {

	// recherche du code d'authorization
	authorize_code, err := gorm.G[models.AuthorizationCode](store.db).Preload("Session.User", nil).Preload(clause.Associations, nil).Where(&models.AuthorizationCode{Code: code}).First(ctx)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fosite.ErrNotFound
		}
		return nil, err
	}

	//marshalling du formulaire
	var form url.Values
	err = json.Unmarshal(authorize_code.Form, &form)
	if err != nil {
		return nil, fmt.Errorf("erreur de unmasharlling du formulaire : %w", err)
	}

	//initialisation de fosite Request
	rq := &fosite.Request{
		ID:                authorize_code.ID.String(),
		RequestedAt:       authorize_code.RequestedAt,
		Client:            &authorize_code.Client,
		RequestedScope:    fosite.Arguments(authorize_code.RequestedScopes),
		GrantedScope:      fosite.Arguments(authorize_code.GrantedScopes),
		Form:              form,
		Session:           &authorize_code.Session,
		RequestedAudience: fosite.Arguments(authorize_code.RequestedAudience),
		GrantedAudience:   fosite.Arguments(authorize_code.GrantedAudience),
	}

	//verification de la validité du code d'authorization
	if authorize_code.Active != nil && !*authorize_code.Active {
		return rq, fosite.ErrInvalidatedAuthorizeCode
	}

	return rq, nil
}

func (store *Store) InvalidateAuthorizeCodeSession(ctx context.Context, code string) (err error) {

	authorize_code, err := gorm.G[models.AuthorizationCode](store.db).Preload("Session.User", nil).Preload(clause.Associations, nil).Where(&models.AuthorizationCode{Code: code}).First(ctx)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fosite.ErrNotFound
		}
		return err
	}

	_, err = gorm.G[models.AuthorizationCode](store.db).Where(&models.AuthorizationCode{ID: authorize_code.ID}).Updates(ctx, models.AuthorizationCode{Active: utils.PtrBool(false)})

	if err != nil {
		return fmt.Errorf("erreur d'invalidation du token: %w", err)
	}

	return nil
}

// AccessTokenStorage
func (store *Store) CreateAccessTokenSession(ctx context.Context, signature string, request fosite.Requester) (err error) {
	parsedID, err := uuid.Parse(request.GetID())
	if err != nil {
		return fmt.Errorf("request id invalide: %w", err)
	}
	client := request.GetClient()

	clientID, err := uuid.Parse(client.GetID())
	if err != nil {
		return fmt.Errorf("client id invalide: %w", err)
	}

	form, err := json.Marshal(request.GetRequestForm())
	if err != nil {
		return fmt.Errorf("erreur de marshalling du access_token form: %w", err)
	}

	session := request.GetSession().(*models.Session)
	if err = gorm.G[models.Session](store.db).Create(ctx, session); err != nil {
		return fmt.Errorf("erreur de création de la sesion pour access_token: %w", err)
	}

	data := models.AccessToken{
		ID:                parsedID,
		Active:            utils.PtrBool(true),
		Signature:         signature,
		RequestedAt:       request.GetRequestedAt().UTC(),
		ClientID:          clientID,
		RequestedScopes:   pq.StringArray(request.GetRequestedScopes()),
		GrantedScopes:     pq.StringArray(request.GetGrantedScopes()),
		Form:              form,
		SessionID:         &session.ID,
		RequestedAudience: pq.StringArray(request.GetRequestedAudience()),
		GrantedAudience:   pq.StringArray(request.GetGrantedAudience()),
	}

	if err = gorm.G[models.AccessToken](store.db).Create(ctx, &data); err != nil {
		return fmt.Errorf("erreur de création d'access_token: %w", err)
	}
	return nil
}

func (store *Store) GetAccessTokenSession(ctx context.Context, signature string, session fosite.Session) (request fosite.Requester, err error) {

	access_token, err := gorm.G[models.AccessToken](store.db).Preload("Session.User", nil).Preload(clause.Associations, nil).Where(&models.AccessToken{Signature: signature}).First(ctx)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fosite.ErrNotFound
		}
		return nil, err
	}

	var form url.Values
	err = json.Unmarshal(access_token.Form, &form)
	if err != nil {
		return nil, fmt.Errorf("erreur de unmasharlling du formulaire : %w", err)
	}

	rq := &fosite.Request{
		ID:                access_token.ID.String(),
		RequestedAt:       access_token.RequestedAt,
		Client:            &access_token.Client,
		RequestedScope:    fosite.Arguments(access_token.RequestedScopes),
		GrantedScope:      fosite.Arguments(access_token.GrantedScopes),
		Form:              form,
		Session:           &access_token.Session,
		RequestedAudience: fosite.Arguments(access_token.RequestedAudience),
		GrantedAudience:   fosite.Arguments(access_token.GrantedAudience),
	}

	return rq, nil
}

func (store *Store) DeleteAccessTokenSession(ctx context.Context, signature string) (err error) {
	if _, err := gorm.G[models.AccessToken](store.db.Unscoped()).Where(&models.AccessToken{Signature: signature}).Delete(ctx); err != nil {
		return fmt.Errorf("erreur de supression de l'accessToken: %w", err)
	}

	return nil
}

// RefreshTokenStorage
func (store *Store) CreateRefreshTokenSession(ctx context.Context, signature string, accessSignature string, request fosite.Requester) (err error) {
	parsedID, err := uuid.Parse(request.GetID())
	if err != nil {
		return fmt.Errorf("request id invalide: %w", err)
	}
	client := request.GetClient()

	clientID, err := uuid.Parse(client.GetID())
	if err != nil {
		return fmt.Errorf("client id invalide: %w", err)
	}

	form, err := json.Marshal(request.GetRequestForm())
	if err != nil {
		return fmt.Errorf("erreur de marshalling du refresh token form: %w", err)
	}

	session := request.GetSession().(*models.Session)
	if err = gorm.G[models.Session](store.db).Create(ctx, session); err != nil {
		return fmt.Errorf("erreur de création de la session: %w", err)
	}

	data := models.RefreshToken{
		ID:                parsedID,
		Active:            utils.PtrBool(true),
		Signature:         signature,
		AccessSignature:   accessSignature,
		RequestedAt:       request.GetRequestedAt().UTC(),
		ClientID:          clientID,
		RequestedScopes:   pq.StringArray(request.GetRequestedScopes()),
		GrantedScopes:     pq.StringArray(request.GetGrantedScopes()),
		Form:              form,
		SessionID:         &session.ID,
		RequestedAudience: pq.StringArray(request.GetRequestedAudience()),
		GrantedAudience:   pq.StringArray(request.GetGrantedAudience()),
	}

	if err = gorm.G[models.RefreshToken](store.db).Create(ctx, &data); err != nil {
		return fmt.Errorf("erreur de cration du refresh token : %w", err)
	}

	return nil
}

func (store *Store) GetRefreshTokenSession(ctx context.Context, signature string, session fosite.Session) (request fosite.Requester, err error) {

	result, err := gorm.G[models.RefreshToken](store.db).Preload("Session.User", nil).Preload(clause.Associations, nil).Where(&models.RefreshToken{Signature: signature}).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fosite.ErrNotFound
		}
		return nil, err
	}

	var form url.Values
	err = json.Unmarshal(result.Form, &form)
	if err != nil {
		return nil, fmt.Errorf("erreur de marshalling du refresk token form: %w", err)
	}

	rq := &fosite.Request{
		ID:                result.ID.String(),
		RequestedAt:       result.RequestedAt,
		Client:            &result.Client,
		RequestedScope:    fosite.Arguments(result.RequestedScopes),
		GrantedScope:      fosite.Arguments(result.GrantedScopes),
		Form:              form,
		Session:           &result.Session,
		RequestedAudience: fosite.Arguments(result.RequestedAudience),
		GrantedAudience:   fosite.Arguments(result.GrantedAudience),
	}

	if result.Active != nil && !*result.Active {
		return rq, fosite.ErrInactiveToken
	}

	return rq, nil
}

func (store *Store) DeleteRefreshTokenSession(ctx context.Context, signature string) (err error) {
	if _, err := gorm.G[models.RefreshToken](store.db.Unscoped()).Where(&models.RefreshToken{Signature: signature}).Delete(ctx); err != nil {
		return fmt.Errorf("erreur de refresh token : %w", err)
	}

	return nil
}

func (store *Store) RotateRefreshToken(ctx context.Context, requestID string, refreshTokenSignature string) (err error) {
	refreshToken, err := gorm.G[models.RefreshToken](store.db).Where(&models.RefreshToken{Signature: refreshTokenSignature}).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fosite.ErrNotFound
		}
		return fmt.Errorf("erreur de recherche du refresh token: %w", err)
	}

	_, err = gorm.G[models.RefreshToken](store.db).Where(&models.RefreshToken{ID: refreshToken.ID}).Updates(ctx, models.RefreshToken{Active: utils.PtrBool(false)})
	if err != nil {
		return fmt.Errorf("erreur d'invalidation du refresh token: %w", err)
	}

	return nil
}

func (store *Store) Authenticate(ctx context.Context, name string, secret string) error {

	user, err := gorm.G[models.User](store.db).Where("user_name = ?", name).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fosite.ErrNotFound.WithDebug("Invalid credentials")
		}
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(secret)); err != nil {
		return fosite.ErrNotFound.WithDebug("Invalid credentials")
	}

	return nil

}

func (store *Store) GetUser(ctx context.Context, username string) (*models.User, error) {
	results, err := gorm.G[models.User](store.db).Preload("Roles.Scopes", nil).Preload(clause.Associations, nil).Where("user_name = ?", username).First(ctx)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fosite.ErrNotFound
		}
		return nil, err
	}

	return &results, err
}
