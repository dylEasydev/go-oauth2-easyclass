package db

//package db

import (
	"context"
	"errors"
	"time"

	"github.com/dylEasydev/go-oauth2-easyclass/db/models"
	"github.com/dylEasydev/go-oauth2-easyclass/utils"
	"github.com/go-jose/go-jose/v3"
	"github.com/ory/fosite"
	"gorm.io/gorm"
)

// GetPublicKey retourne la clé publique pour un issuer, subject et keyId spécifique
func (store *Store) GetPublicKey(ctx context.Context, issuer, subject, keyId string) (*jose.JSONWebKey, error) {

	key, err := gorm.G[models.ClientKey](store.db).Where(&models.ClientKey{Issuer: issuer, Subject: subject, KeyID: keyId}).First(ctx)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fosite.ErrNotFound
		}
		return nil, err
	}

	return (*jose.JSONWebKey)(&key.JWK), nil
}

// GetPublicKeys retourne toutes les clés publiques pour un issuer et subject
func (store *Store) GetPublicKeys(ctx context.Context, issuer, subject string) (*jose.JSONWebKeySet, error) {

	keys, err := gorm.G[models.ClientKey](store.db).Where(&models.ClientKey{Issuer: issuer, Subject: subject}).Find(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fosite.ErrNotFound
		}
		return nil, err
	}

	jwks := []jose.JSONWebKey{}
	for _, k := range keys {
		jwks = append(jwks, jose.JSONWebKey(k.JWK))

	}
	return &jose.JSONWebKeySet{Keys: jwks}, nil
}

// GetPublicKeyScopes retourne les scopes assignés à une clé publique
func (store *Store) GetPublicKeyScopes(ctx context.Context, issuer, subject, keyId string) ([]string, error) {

	key, err := gorm.G[models.ClientKey](store.db).Where(&models.ClientKey{Issuer: issuer, Subject: subject, KeyID: keyId}).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return key.Scopes, nil
}

// IsJWTUsed retourne true si le JWT est déjà utilisé ou expiré
func (store *Store) IsJWTUsed(ctx context.Context, jti string) (bool, error) {
	clientJwt, err := gorm.G[models.ClientJWT](store.db).Where(&models.ClientJWT{JTI: jti}).First(ctx)
	if err != nil {
		return false, err
	}
	// Retourne true si le JTI existe et est encore valide (donc déjà utilisé).
	return !clientJwt.IsValid(), nil
}

// MarkJWTUsedForTime marque un JWT comme utilisable jusqu'à exp
func (store *Store) MarkJWTUsedForTime(ctx context.Context, jti string, exp time.Time) error {
	_, err := gorm.G[models.ClientJWT](store.db).Where(&models.ClientJWT{JTI: jti}).Where("expires_at > ?", time.Now().UTC()).First(ctx)
	if err != nil {
		return fosite.ErrJTIKnown
	}
	jwt := models.ClientJWT{
		JTI:       jti,
		ExpiresAt: exp.UTC(),
		Active:    utils.PtrBool(true),
	}

	return gorm.G[models.ClientJWT](store.db).Create(ctx, &jwt)
}
