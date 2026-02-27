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
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (store *Store) CreatePKCERequestSession(ctx context.Context, signature string, requester fosite.Requester) error {
	parsedID, err := uuid.Parse(requester.GetID())
	if err != nil {
		return fmt.Errorf("request id invalide: %w", err)
	}
	client := requester.GetClient()

	clientID, err := uuid.Parse(client.GetID())
	if err != nil {
		return fmt.Errorf("client id invalide: %w", err)
	}

	formValues := requester.GetRequestForm()
	form, err := json.Marshal(formValues)
	if err != nil {
		return fmt.Errorf("errreur de marchalling du pcke form %w", err)
	}

	session := requester.GetSession().(*models.Session)
	if err = gorm.G[models.Session](store.db).Create(ctx, session); err != nil {
		return fmt.Errorf("erreur de cration de la session: %w", err)
	}

	data := models.PKCE{
		ID:                parsedID,
		Active:            utils.PtrBool(true),
		Signature:         signature,
		RequestedAt:       requester.GetRequestedAt().UTC(),
		ClientID:          clientID,
		RequestedScopes:   pq.StringArray(requester.GetRequestedScopes()),
		GrantedScopes:     pq.StringArray(requester.GetGrantedScopes()),
		Form:              form,
		SessionID:         &session.ID,
		RequestedAudience: pq.StringArray(requester.GetRequestedAudience()),
		GrantedAudience:   pq.StringArray(requester.GetGrantedAudience()),
	}

	if err = gorm.G[models.PKCE](store.db).Create(ctx, &data); err != nil {
		return fmt.Errorf("erreur de creation du PCKE: %w", err)
	}

	return nil
}

func (store *Store) GetPKCERequestSession(ctx context.Context, signature string, session fosite.Session) (fosite.Requester, error) {
	var result models.PKCE

	result, err := gorm.G[models.PKCE](store.db).Joins(clause.JoinTarget{Association: "Client"}, nil).Joins(clause.JoinTarget{Association: "Session"}, nil).Joins(clause.JoinTarget{Association: "Session.User"}, nil).Where(&models.PKCE{Signature: signature}).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fosite.ErrNotFound
		}
		return nil, err
	}

	var form url.Values
	err = json.Unmarshal(result.Form, &form)
	if err != nil {
		return nil, fmt.Errorf("erreur unmarshal du formulaire PCKE: %w", err)
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

	return rq, nil
}

func (store *Store) DeletePKCERequestSession(ctx context.Context, signature string) error {
	if _, err := gorm.G[models.PKCE](store.db.Unscoped()).Where(&models.PKCE{Signature: signature}).Delete(ctx); err != nil {
		return fmt.Errorf("erreur de suppression de PCKE request: %w", err)
	}

	return nil
}
