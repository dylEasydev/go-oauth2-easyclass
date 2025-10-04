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
	"github.com/ory/fosite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

//implementation de PARStorage

func (store *Store) CreatePARSession(ctx context.Context, requestURI string, request fosite.AuthorizeRequester) error {
	client := request.GetClient()

	form, err := json.Marshal(request.GetRequestForm())
	if err != nil {
		return fmt.Errorf("erreur de marshalling du PAR form : %w", err)
	}

	session := request.GetSession().(*models.Session)

	if err = store.db.WithContext(ctx).
		Create(session).Error; err != nil {
		return fmt.Errorf("erreur de création de la sesion pour PAR: %w", err)
	}

	data := models.PARRequest{
		ID:         uuid.MustParse(request.GetID()),
		RequestURI: requestURI,
		Form:       form,
		ClientID:   uuid.MustParse(client.GetID()),
		SessionID:  session.ID,
	}

	if err = store.db.WithContext(ctx).
		Create(&data).Error; err != nil {
		return fmt.Errorf("erreur de création de PAR request: %w", err)
	}

	return nil
}

func (store *Store) GetPARSession(ctx context.Context, requestURI string) (fosite.AuthorizeRequester, error) {
	var par models.PARRequest

	if err := store.db.WithContext(ctx).
		Preload("Session.User").
		Preload(clause.Associations).
		Where(&models.PARRequest{RequestURI: requestURI}).
		First(&par).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fosite.ErrNotFound
		}
		return nil, err
	}

	var form url.Values
	err := json.Unmarshal(par.Form, &form)
	if err != nil {
		return nil, fmt.Errorf("erreur de unmasharlling du formulaire : %w", err)
	}

	rq := &fosite.AuthorizeRequest{
		Request: fosite.Request{
			ID:          par.ID.String(),
			RequestedAt: time.Now(),
			Client:      &par.Client,
			Form:        form,
			Session:     &par.Session,
		},
	}

	if par.Used {
		return rq, fosite.ErrInvalidRequest.WithHint("ce PAR request est déjà utilisé")
	}
	if time.Now().After(par.ExpiresAt) {
		return nil, fosite.ErrInvalidRequest.WithHint("ce PAR request est expiré.")
	}

	return rq, nil
}
func (store *Store) DeletePARSession(ctx context.Context, requestURI string) (err error) {
	if err := store.db.WithContext(ctx).
		Where(&models.PARRequest{RequestURI: requestURI}).
		Delete(&models.PARRequest{}).Error; err != nil {
		return fmt.Errorf("erreur de supression du PAR: %w", err)
	}

	return nil
}
