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

var (
	ErrNonceExpired = errors.New("nonce expiré")
)

//implementation de l'interface nonceManager
//pour l'extensions  verifiable(Nonce)

// création du nonce
func (store *Store) NewNonce(ctx context.Context, accessToken string, expiresAt time.Time) (string, error) {
	//initialisation de la struture
	data := &models.Nonce{
		AccessToken: accessToken,
		ExpiresAt:   expiresAt.UTC(),
		Nonce:       uuid.New().String(),
	}

	// création du Nonce en BD
	// utilisation des méthodes génériques
	if err := gorm.G[models.Nonce](store.db).Create(ctx, data); err != nil {
		return "", fmt.Errorf("erreur de création du nonce: %w", err)
	}

	return data.Nonce, nil
}

// verifie si access_token correspond au nonce donnée
// et si le nonce n'est pas expiré
func (store *Store) IsNonceValid(ctx context.Context, accessToken string, nonce string) error {

	return store.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// recherche du nonce correspondant à access_token
		result, err := gorm.G[models.Nonce](tx).Where(&models.Nonce{AccessToken: accessToken, Nonce: nonce}).First(ctx)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fosite.ErrNotFound
			}
			return fmt.Errorf("erreur lecture nonce: %w", err)
		}

		// verification de la validité du Nonce (si ExpiresAt est avant maintenant => expiré)
		if result.ExpiresAt.Before(time.Now().UTC()) {
			// suppression du nonce expiré
			_, _ = gorm.G[models.Nonce](tx.Unscoped()).Where(&models.Nonce{ID: result.ID}).Delete(ctx)
			return ErrNonceExpired
		}
		if _, err := gorm.G[models.Nonce](tx.Unscoped()).Where(&models.Nonce{ID: result.ID}).Delete(ctx); err != nil {
			return fmt.Errorf("erreur suppression nonce: %w", err)
		}

		return nil
	})
}
