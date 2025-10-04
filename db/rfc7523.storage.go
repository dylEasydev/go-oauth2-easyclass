package db

import (
	"context"
	"errors"
	"time"

	"github.com/dylEasydev/go-oauth2-easyclass/db/models"
	"github.com/dylEasydev/go-oauth2-easyclass/utils"
	"github.com/go-jose/go-jose/v3"
	"gorm.io/gorm"
)

// GetPublicKey retourne la clé publique pour un issuer, subject et keyId spécifique
func (store *Store) GetPublicKey(ctx context.Context, issuer, subject, keyId string) (*jose.JSONWebKey, error) {
	var key models.ClientKey
	err := store.db.WithContext(ctx).
		Where(&models.ClientKey{Issuer: issuer, Subject: subject, KeyID: keyId}).
		First(&key).Error
	if err != nil {
		return nil, err
	}

	var jwk jose.JSONWebKey
	if err := jwk.UnmarshalJSON(key.JWK); err != nil {
		return nil, err
	}
	return &jwk, nil
}

// GetPublicKeys retourne toutes les clés publiques pour un issuer et subject
func (store *Store) GetPublicKeys(ctx context.Context, issuer, subject string) (*jose.JSONWebKeySet, error) {
	var keys []models.ClientKey
	err := store.db.WithContext(ctx).
		Where(&models.ClientKey{Issuer: issuer, Subject: subject}).
		Find(&keys).Error
	if err != nil {
		return nil, err
	}

	jwks := []jose.JSONWebKey{}
	for _, k := range keys {
		var jwk jose.JSONWebKey
		if err := jwk.UnmarshalJSON(k.JWK); err == nil {
			jwks = append(jwks, jwk)
		}
	}
	return &jose.JSONWebKeySet{Keys: jwks}, nil
}

// GetPublicKeyScopes retourne les scopes assignés à une clé publique
func (store *Store) GetPublicKeyScopes(ctx context.Context, issuer, subject, keyId string) ([]string, error) {
	var key models.ClientKey
	err := store.db.WithContext(ctx).
		Where(&models.ClientKey{Issuer: issuer, Subject: subject, KeyID: keyId}).
		First(&key).Error
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
	var jwt models.ClientJWT
	err := store.db.WithContext(ctx).Where(&models.ClientJWT{JTI: jti}).First(&jwt).Error
	if err != nil {
		return false, err
	}
	return !jwt.IsValid(), nil
}

// MarkJWTUsedForTime marque un JWT comme utilisé jusqu'à exp
func (s *Store) MarkJWTUsedForTime(ctx context.Context, jti string, exp time.Time) error {
	jwt := models.ClientJWT{
		JTI:       jti,
		ExpiresAt: exp,
		Active:    utils.PtrBool(true),
	}

	return s.db.WithContext(ctx).Create(&jwt).Error
}
