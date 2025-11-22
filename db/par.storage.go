package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/dylEasydev/go-oauth2-easyclass/db/models"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/ory/fosite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

//implementation de PARStorage

func (store *Store) CreatePARSession(ctx context.Context, requestURI string, request fosite.AuthorizeRequester) error {
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
		return fmt.Errorf("erreur de marshalling du PAR form : %w", err)
	}

	redirectUri, err := json.Marshal(request.GetRedirectURI())
	if err != nil {
		return fmt.Errorf("erreur de marshalling du PAR redirect URI : %w", err)
	}

	session := request.GetSession().(*models.Session)

	if err = gorm.G[models.Session](store.db).Create(ctx, session); err != nil {
		return fmt.Errorf("erreur de création de la sesion pour PAR: %w", err)
	}

	data := models.PARRequest{
		ID:                parsedID,
		RequestURI:        requestURI,
		RequestedAt:       request.GetRequestedAt().UTC(),
		Form:              form,
		RequestedScopes:   pq.StringArray(request.GetRequestedScopes()),
		GrantedScopes:     pq.StringArray(request.GetGrantedScopes()),
		ClientID:          clientID,
		SessionID:         session.ID,
		RequestedAudience: pq.StringArray(request.GetRequestedAudience()),
		GrantedAudience:   pq.StringArray(request.GetGrantedAudience()),
		RedirectURI:       redirectUri,
		ResponseMode:      string(request.GetResponseMode()),
	}

	if err = gorm.G[models.PARRequest](store.db).Create(ctx, &data); err != nil {
		return fmt.Errorf("erreur de création de PAR request: %w", err)
	}

	return nil
}

func (store *Store) GetPARSession(ctx context.Context, requestURI string) (fosite.AuthorizeRequester, error) {

	par, err := gorm.G[models.PARRequest](store.db).Preload("Session.User", nil).Preload(clause.Associations, nil).Where(&models.PARRequest{RequestURI: requestURI}).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fosite.ErrNotFound
		}
		return nil, err
	}

	var form url.Values
	var redirectUri url.URL
	err = json.Unmarshal(par.Form, &form)
	if err != nil {
		return nil, fmt.Errorf("erreur de unmasharlling du formulaire : %w", err)
	}
	err = json.Unmarshal(par.RedirectURI, &redirectUri)
	if err != nil {
		return nil, fmt.Errorf("erreur de unmasharlling du redirectURI : %w", err)
	}
	rq := &fosite.AuthorizeRequest{
		Request: fosite.Request{
			ID:                par.ID.String(),
			RequestedAt:       par.RequestedAt,
			Client:            &par.Client,
			RequestedScope:    fosite.Arguments(par.RequestedScopes),
			GrantedScope:      fosite.Arguments(par.GrantedScopes),
			Form:              form,
			Session:           &par.Session,
			RequestedAudience: fosite.Arguments(par.RequestedAudience),
			GrantedAudience:   fosite.Arguments(par.GrantedAudience),
		},
		ResponseTypes: par.Client.GetResponseTypes(),
		RedirectURI:   &redirectUri,
		ResponseMode:  fosite.ResponseModeType(par.ResponseMode),
	}

	if par.Used {
		return rq, fosite.ErrInvalidRequest.WithHint("ce PAR request est déjà utilisé")
	}
	if time.Now().UTC().After(par.ExpiresAt) {
		return nil, fosite.ErrInvalidRequest.WithHint("ce PAR request est expiré.")
	}

	return rq, nil
}
func (store *Store) DeletePARSession(ctx context.Context, requestURI string) (err error) {
	if _, err := gorm.G[models.PARRequest](store.db.Unscoped()).Where(&models.PARRequest{RequestURI: requestURI}).Delete(ctx); err != nil {
		return fmt.Errorf("erreur de supression du PAR: %w", err)
	}

	return nil
}
