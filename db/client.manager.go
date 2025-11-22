package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dylEasydev/go-oauth2-easyclass/db/models"
	"github.com/google/uuid"
	"github.com/ory/fosite"
	"gorm.io/gorm"
)

//client manager

func (store *Store) GetClient(ctx context.Context, id string) (fosite.Client, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("client id invalide: %w", err)
	}

	client, err := gorm.G[models.Client](store.db).Preload("Keys", nil).Where(&models.Client{ID: parsedID}).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fosite.ErrNotFound
		}
		return nil, err
	}
	return &client, nil
}

func (store *Store) ClientAssertionJWTValid(ctx context.Context, jti string) error {
	clientJwt, err := gorm.G[models.ClientJWT](store.db).Where(&models.ClientJWT{JTI: jti}).First(ctx)

	if err != nil {
		return nil
		/*if errors.Is(err, gorm.ErrRecordNotFound) {
			return fosite.ErrNotFound
		}
		return err*/
	}
	// Si le JTI existe et est encore valide, il a déjà été utilisé : rejeter.
	if clientJwt.IsValid() {
		return fosite.ErrJTIKnown
	}
	return nil
}

func (store *Store) SetClientAssertionJWT(ctx context.Context, jti string, exp time.Time) error {
	return store.MarkJWTUsedForTime(ctx, jti, exp)
}
