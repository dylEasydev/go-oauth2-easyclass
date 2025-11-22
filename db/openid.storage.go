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

	// recherche du code d'authorisation correspondant (authorizeCode)
	authorizeCodeModel, err := gorm.G[models.AuthorizationCode](store.db).Preload("Session.User", nil).Preload(clause.Associations, nil).Where(&models.AuthorizationCode{Code: authorizeCode}).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fosite.ErrNotFound
		}
		return nil, err
	}

	//unmarshal du formulaire
	var form url.Values
	if err := json.Unmarshal(authorizeCodeModel.Form, &form); err != nil {
		return nil, fmt.Errorf("erreur de unmarshal du formulaire openid : %w", err)
	}

	//intialisation du fosite Request
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

	//verification de la validité du code d'authorization
	if authorizeCodeModel.Active != nil && !*authorizeCodeModel.Active {
		return rq, fosite.ErrInvalidatedAuthorizeCode
	}

	return rq, nil
}

// suppression du code de la session openid
func (store *Store) DeleteOpenIDConnectSession(ctx context.Context, authorizeCode string) error {

	//suppression du code d'authorization de la BD
	if _, err := gorm.G[models.AuthorizationCode](store.db.Unscoped()).Where(&models.AuthorizationCode{Code: authorizeCode}).Delete(ctx); err != nil {
		return fmt.Errorf("erreur de supression du code d'authorization: %w", err)
	}

	return nil
}
