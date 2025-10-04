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

//implementation de l'inreface de CoreStorage

// AuthorizationStorage
func (store *Store) CreateAuthorizeCodeSession(ctx context.Context, code string, request fosite.Requester) (err error) {
	client := request.GetClient()

	form, err := json.Marshal(request.GetRequestForm())
	if err != nil {
		return fmt.Errorf("erreur de marshalling du authorize form : %w", err)
	}

	session := request.GetSession().(*models.Session)

	if err = store.db.WithContext(ctx).
		Create(session).Error; err != nil {
		return fmt.Errorf("erreur de création de la sesion pour authorize code: %w", err)
	}

	data := models.AuthorizationCode{
		ID:                uuid.MustParse(request.GetID()),
		Active:            utils.PtrBool(true),
		Code:              code,
		RequestedAt:       request.GetRequestedAt(),
		ClientID:          uuid.MustParse(client.GetID()),
		RequestedScopes:   pq.StringArray(request.GetRequestedScopes()),
		GrantedScopes:     pq.StringArray(request.GetGrantedScopes()),
		Form:              form,
		SessionID:         session.ID,
		RequestedAudience: pq.StringArray(request.GetRequestedAudience()),
		GrantedAudience:   pq.StringArray(request.GetGrantedAudience()),
	}

	if err = store.db.WithContext(ctx).
		Create(&data).Error; err != nil {
		return fmt.Errorf("erreur de création d'authorize code: %w", err)
	}

	return nil
}

func (store *Store) GetAuthorizeCodeSession(ctx context.Context, code string, session fosite.Session) (request fosite.Requester, err error) {

	var authorize_code models.AuthorizationCode

	if err := store.db.WithContext(ctx).
		Preload("Session.User").
		Preload(clause.Associations).
		Where(&models.AuthorizationCode{Code: code}).
		First(&authorize_code).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fosite.ErrNotFound
		}
		return nil, err
	}

	var form url.Values
	err = json.Unmarshal(authorize_code.Form, &form)
	if err != nil {
		return nil, fmt.Errorf("erreur de unmasharlling du formulaire : %w", err)
	}

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

	if authorize_code.Active != nil && !*authorize_code.Active {
		return rq, fosite.ErrInvalidatedAuthorizeCode
	}

	return rq, nil
}

func (store *Store) InvalidateAuthorizeCodeSession(ctx context.Context, code string) (err error) {
	var authorize_code models.AuthorizationCode

	if err := store.db.WithContext(ctx).
		Where(&models.AuthorizationCode{Code: code}).
		First(&authorize_code).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fosite.ErrNotFound
		}
		return err
	}

	authorize_code.Active = utils.PtrBool(false)
	if err := store.db.WithContext(ctx).
		Save(authorize_code).Error; err != nil {
		return fmt.Errorf("erreur d'invalidation du token: %w", err)
	}

	return nil
}

// AccessTokenStorage
func (store *Store) CreateAccessTokenSession(ctx context.Context, signature string, request fosite.Requester) (err error) {
	client := request.GetClient()

	form, err := json.Marshal(request.GetRequestForm())
	if err != nil {
		return fmt.Errorf("erreur de marshalling du access_token form: %w", err)
	}

	session := request.GetSession().(*models.Session)
	if err = store.db.WithContext(ctx).
		Create(&session).Error; err != nil {
		return fmt.Errorf("erreur de création de la sesion pour access_token: %w", err)
	}

	data := models.AccessToken{
		ID:                uuid.MustParse(request.GetID()),
		Active:            utils.PtrBool(true),
		Signature:         signature,
		RequestedAt:       request.GetRequestedAt(),
		ClientID:          uuid.MustParse(client.GetID()),
		RequestedScopes:   pq.StringArray(request.GetRequestedScopes()),
		GrantedScopes:     pq.StringArray(request.GetGrantedScopes()),
		Form:              form,
		SessionID:         &session.ID,
		RequestedAudience: pq.StringArray(request.GetRequestedAudience()),
		GrantedAudience:   pq.StringArray(request.GetGrantedAudience()),
	}

	if err = store.db.WithContext(ctx).
		Create(&data).Error; err != nil {
		return fmt.Errorf("erreur de création d'access_token: %w", err)
	}
	return nil
}

func (store *Store) GetAccessTokenSession(ctx context.Context, signature string, session fosite.Session) (request fosite.Requester, err error) {

	var access_token models.AccessToken

	if err := store.db.WithContext(ctx).
		Preload("Session.User").
		Preload(clause.Associations).
		Where(&models.AccessToken{Signature: signature}).
		First(&access_token).Error; err != nil {
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
	if err := store.db.WithContext(ctx).
		Where(&models.AccessToken{Signature: signature}).
		Delete(&models.AccessToken{}).Error; err != nil {
		return fmt.Errorf("erreur de supression de l'accessToken: %w", err)
	}

	return nil
}

// RefreshTokenStorage
func (store *Store) CreateRefreshTokenSession(ctx context.Context, signature string, request fosite.Requester) (err error) {
	client := request.GetClient()

	form, err := json.Marshal(request.GetRequestForm())
	if err != nil {
		return fmt.Errorf("erreur de marshalling du refresh token form: %w", err)
	}

	session := request.GetSession().(*models.Session)
	if err = store.db.WithContext(ctx).
		Create(&session).Error; err != nil {
		return fmt.Errorf("erreur de création de la session: %w", err)
	}

	data := models.RefreshToken{
		ID:                uuid.MustParse(request.GetID()),
		Active:            utils.PtrBool(true),
		Signature:         signature,
		RequestedAt:       request.GetRequestedAt(),
		ClientID:          uuid.MustParse(client.GetID()),
		RequestedScopes:   pq.StringArray(request.GetRequestedScopes()),
		GrantedScopes:     pq.StringArray(request.GetGrantedScopes()),
		Form:              form,
		SessionID:         &session.ID,
		RequestedAudience: pq.StringArray(request.GetRequestedAudience()),
		GrantedAudience:   pq.StringArray(request.GetGrantedAudience()),
	}

	if err = store.db.WithContext(ctx).
		Create(&data).Error; err != nil {
		return fmt.Errorf("erreur de cration du refresh token : %w", err)
	}

	return nil
}

func (store *Store) GetRefreshTokenSession(ctx context.Context, signature string, session fosite.Session) (request fosite.Requester, err error) {
	var result models.RefreshToken

	if err := store.db.WithContext(ctx).
		Preload("Session.User").
		Preload(clause.Associations).
		Where(&models.RefreshToken{Signature: signature}).
		First(&result).Error; err != nil {
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
	if err := store.db.WithContext(ctx).
		Where(&models.RefreshToken{Signature: signature}).
		Delete(&models.RefreshToken{}).Error; err != nil {
		return fmt.Errorf("erreur de refresh token : %w", err)
	}

	return nil
}

func (store *Store) Authenticate(ctx context.Context, name string, secret string) error {
	var user models.User
	if err := store.db.WithContext(ctx).Where("username = ?", name).First(&user).Error; err != nil {
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
