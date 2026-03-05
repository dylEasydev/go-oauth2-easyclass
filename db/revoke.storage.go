package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/dylEasydev/go-oauth2-easyclass/db/models"
	"github.com/dylEasydev/go-oauth2-easyclass/utils"
	"github.com/ory/fosite"
	"gorm.io/gorm"
)

func (store *Store) RevokeRefreshToken(ctx context.Context, requestID string) error {

	if err := store.db.WithContext(ctx).Model(&models.RefreshToken{}).Where(&models.RefreshToken{RequestId: requestID}).Updates(&models.RefreshToken{Active: utils.PtrBool(false)}).Error; err != nil {
		return fmt.Errorf("erreur de revocation du jeton de rafraichissement : %w", err)
	}
	return nil
}

func (store *Store) RevokeRefreshTokenMaybeGracePeriod(ctx context.Context, requestID string, signature string) error {
	refreshToken, err := gorm.G[models.RefreshToken](store.db).Where(&models.RefreshToken{RequestId: requestID, Signature: signature}).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fosite.ErrNotFound
		}
		return fmt.Errorf("erreur de revocation du jeton de rafraichissement : %w", err)
	}

	refreshToken.Active = utils.PtrBool(false)
	if err := store.db.WithContext(ctx).Model(&models.RefreshToken{}).Save(refreshToken).Error; err != nil {
		return fmt.Errorf("erreur de revocation du jeton de rafraichissement : %w", err)
	}

	return nil
}

func (store *Store) RevokeAccessToken(ctx context.Context, requestID string) error {
	if err := store.db.WithContext(ctx).Model(&models.AccessToken{}).Where(&models.AccessToken{RequestId: requestID}).Updates(&models.AccessToken{Active: utils.PtrBool(false)}).Error; err != nil {
		return fmt.Errorf("erreur de revocation du jeton de rafraichissement : %w", err)
	}
	return nil
}
