package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/dylEasydev/go-oauth2-easyclass/db/models"
	"github.com/dylEasydev/go-oauth2-easyclass/utils"
	"github.com/google/uuid"
	"github.com/ory/fosite"
	"gorm.io/gorm"
)

func (store *Store) RevokeRefreshToken(ctx context.Context, requestID string) error {

	var result models.RefreshToken

	id := uuid.MustParse(requestID)
	if err := store.db.WithContext(ctx).
		Where(&models.RefreshToken{ID: id}).
		First(&result).Error; err != nil {
		return fmt.Errorf("%w: %w", fosite.ErrNotFound, err)
	}

	result.Active = utils.PtrBool(false)
	if err := store.db.WithContext(ctx).
		Save(result).Error; err != nil {
		return fmt.Errorf("erreur inalidation du refresh token: %w", err)
	}

	return nil
}

func (store *Store) RevokeRefreshTokenMaybeGracePeriod(ctx context.Context, requestID string, signature string) error {
	return store.RevokeRefreshToken(ctx, requestID)
}

func (store *Store) RevokeAccessToken(ctx context.Context, requestID string) error {
	var result models.AccessToken

	id := uuid.MustParse(requestID)
	if err := store.db.WithContext(ctx).
		Where(&models.AccessToken{ID: id}).
		First(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fosite.ErrNotFound
		}
		return err
	}

	result.Active = utils.PtrBool(false)
	if err := store.db.WithContext(ctx).
		Save(result).Error; err != nil {
		return fmt.Errorf("erreur d'invalidation du jeton: %w", err)
	}

	return nil
}
