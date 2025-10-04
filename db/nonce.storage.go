package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dylEasydev/go-oauth2-easyclass/db/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

//implementation de l'interface nonceManager
//pour l'extensions  verifiable(Nonce)

// création du nonce
func (store *Store) NewNonce(ctx context.Context, accessToken string, expiresAt time.Time) (string, error) {
	data := &models.Nonce{
		AccessToken: accessToken,
		ExpiresAt:   expiresAt,
		Nonce:       uuid.New().String(),
	}

	if err := store.db.WithContext(ctx).Create(data).Error; err != nil {
		return "", fmt.Errorf("erreur de creation du Nonce")
	}

	return data.Nonce, nil
}

// verifie si access_token correspond au nonce donnée
// et si le nonce n'est pas expiré
func (store *Store) IsNonceValid(ctx context.Context, accessToken string, nonce string) error {
	var result models.Nonce

	if err := store.db.WithContext(ctx).
		Where(&models.Nonce{AccessToken: accessToken, Nonce: nonce}).
		First(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("le nonce fourni n'est pas valide ")
		}
		return err
	}

	if result.ExpiresAt.After(time.Now()) {
		return fmt.Errorf("le nonce est expiré")
	}

	return nil
}
