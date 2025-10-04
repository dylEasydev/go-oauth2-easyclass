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
	client := requester.GetClient()

	form, err := json.Marshal(requester.GetRequestForm())
	if err != nil {
		return fmt.Errorf("errreur de marchalling du pcke form %w", err)
	}

	session := requester.GetSession().(*models.Session)
	if err = store.db.WithContext(ctx).
		Create(&session).Error; err != nil {
		return fmt.Errorf("erreur de cration de la session: %w", err)
	}

	data := models.PKCE{
		ID:                uuid.MustParse(requester.GetID()),
		Active:            utils.PtrBool(true),
		Signature:         signature,
		RequestedAt:       requester.GetRequestedAt(),
		ClientID:          uuid.MustParse(client.GetID()),
		RequestedScopes:   pq.StringArray(requester.GetRequestedScopes()),
		GrantedScopes:     pq.StringArray(requester.GetGrantedScopes()),
		Form:              form,
		SessionID:         &session.ID,
		RequestedAudience: pq.StringArray(requester.GetRequestedAudience()),
		GrantedAudience:   pq.StringArray(requester.GetGrantedAudience()),
	}

	if err = store.db.WithContext(ctx).
		Create(&data).Error; err != nil {
		return fmt.Errorf("erreur de creation du PCKE: %w", err)
	}

	return nil
}

func (store *Store) GetPKCERequestSession(ctx context.Context, signature string, session fosite.Session) (fosite.Requester, error) {
	var result models.PKCE

	if err := store.db.WithContext(ctx).
		Preload("Session.User").
		Preload(clause.Associations).
		Where(&models.PKCE{Signature: signature}).
		First(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fosite.ErrNotFound
		}
		return nil, err
	}

	var form url.Values
	err := json.Unmarshal(result.Form, &form)
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
	if err := store.db.WithContext(ctx).
		Where(&models.PKCE{Signature: signature}).
		Delete(&models.PKCE{}).Error; err != nil {
		return fmt.Errorf("erreur de suppression de PCKE request: %w", err)
	}

	return nil
}
