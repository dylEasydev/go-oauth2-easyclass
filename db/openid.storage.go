package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	"github.com/dylEasydev/go-oauth2-easyclass/db/models"
	"github.com/ory/fosite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

//implementation de l'interface OpenIDConnectRequestStorage

// création d'une session openid
// identique à la creation d'une session oauth2
func (store *Store) CreateOpenIDConnectSession(ctx context.Context, authorizeCode string, requester fosite.Requester) error {
	return store.CreateAuthorizeCodeSession(ctx, authorizeCode, requester)
}

// recupératon d'une session openid
func (store *Store) GetOpenIDConnectSession(ctx context.Context, authorizeCode string, requester fosite.Requester) (fosite.Requester, error) {
	var authorizeCodeModel models.AuthorizationCode

	if err := store.db.WithContext(ctx).
		Preload("Session.User").
		Preload(clause.Associations).
		Where(&models.AuthorizationCode{Code: authorizeCode}).
		First(&authorizeCodeModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fosite.ErrNotFound
		}
		return nil, err
	}

	var form url.Values
	if err := json.Unmarshal(authorizeCodeModel.Form, &form); err != nil {
		return nil, fmt.Errorf("erreur de unmarshal du formulaire openid : %w", err)
	}

	rq := &fosite.Request{
		ID:                authorizeCodeModel.ID.String(),
		RequestedAt:       authorizeCodeModel.RequestedAt,
		Client:            &authorizeCodeModel.Client,
		RequestedScope:    fosite.Arguments(authorizeCodeModel.RequestedScopes),
		GrantedScope:      fosite.Arguments(authorizeCodeModel.GrantedScopes),
		RequestedAudience: fosite.Arguments(authorizeCodeModel.RequestedAudience),
		GrantedAudience:   fosite.Arguments(authorizeCodeModel.GrantedAudience),
		Form:              form,
		Session:           &authorizeCodeModel.Session,
	}

	if authorizeCodeModel.Active != nil && !*authorizeCodeModel.Active {
		return rq, fosite.ErrInvalidatedAuthorizeCode
	}

	return rq, nil
}

// suppression du code de la session openid
func (store *Store) DeleteOpenIDConnectSession(ctx context.Context, authorizeCode string) error {
	if err := store.db.WithContext(ctx).
		Where(&models.AuthorizationCode{Code: authorizeCode}).
		Delete(&models.AuthorizationCode{}).Error; err != nil {
		return fmt.Errorf("erreur de supression du code d'authorization: %w", err)
	}

	return nil
}
