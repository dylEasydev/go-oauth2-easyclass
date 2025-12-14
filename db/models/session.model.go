package models

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mohae/deepcopy"
	"github.com/ory/fosite"
	"github.com/ory/fosite/token/jwt"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// structure de session basée sur fosite
type Session struct {
	ID uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`

	Username  string
	Subject   string
	ExpiresAt datatypes.JSON `gorm:"type:jsonb"`

	//OIDC specification
	AuthTime time.Time

	AMR datatypes.JSON `gorm:"type:jsonb;default:'[\"pwd\"]'"`
	ACR string         `gorm:"default:'urn:mace:incommon:iap:silver'"`

	Extra datatypes.JSON

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	ClientID uuid.UUID `gorm:"type:uuid;not null"`
	Client   Client    `gorm:"foreignKey:ClientID;references:ID"`

	UserID *uuid.UUID `gorm:"type:uuid"`
	User   User       `gorm:"foreignKey:UserID;references:ID"`
}

func (Session) TableName() string {
	return "sessions"
}

func NewSession(
	ctx context.Context,
	clientID string,
	userID string,
	username string,
	subject string,
	extra map[string]any,
) (*Session, error) {
	idClient, err := uuid.Parse(clientID)
	if err != nil {
		return nil, err
	}
	idUser, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}
	session := &Session{
		UserID:   &idUser,
		ClientID: idClient,
		Username: username,
		Subject:  subject,
		AuthTime: time.Now().UTC(),
	}

	if extra != nil {
		sess_extra, err := json.Marshal(extra)
		if err != nil {
			return nil, fmt.Errorf("error marshalling session extra: %w", err)
		}

		session.Extra = sess_extra
	}

	return session, nil
}

func (s *Session) SetSubject(subject string) {
	s.Subject = subject
}

func (s *Session) SetExpiresAt(key fosite.TokenType, exp time.Time) {
	expiresAt := make(map[fosite.TokenType]time.Time)

	if s.ExpiresAt != nil {
		_ = json.Unmarshal(s.ExpiresAt, &expiresAt)
	}

	expiresAt[key] = exp

	sess_expires, _ := json.Marshal(expiresAt)

	s.ExpiresAt = sess_expires
}

func (s *Session) GetExpiresAt(key fosite.TokenType) time.Time {
	if s.ExpiresAt == nil {
		return time.Time{}
	}

	expiresAt := make(map[fosite.TokenType]time.Time)
	_ = json.Unmarshal(s.ExpiresAt, &expiresAt)

	if _, ok := expiresAt[key]; !ok {
		return time.Time{}
	}

	return expiresAt[key]
}

func (s *Session) GetUsername() string {
	if s == nil {
		return ""
	}

	return s.Username
}

func (s *Session) GetExtraClaims() map[string]interface{} {
	if s == nil {
		return nil
	}

	var extra map[string]interface{}

	if s.Extra != nil {
		err := json.Unmarshal(s.Extra, &extra)
		if err != nil {
			return nil
		}
	}

	return extra
}

func (s *Session) GetSubject() string {
	if s == nil {
		return ""
	}

	return s.Subject
}

func (s *Session) Clone() fosite.Session {
	if s == nil {
		return nil
	}

	return deepcopy.Copy(s).(fosite.Session)
}

func (s *Session) GetJWTClaims() jwt.JWTClaimsContainer {
	if s == nil {
		return &jwt.JWTClaims{}
	}

	claims := jwt.JWTClaims{}

	//subject
	claims.Subject = s.Subject

	//extra
	claims.Extra = s.GetExtraClaims()

	return &claims
}

func (s *Session) IDTokenClaims() *jwt.IDTokenClaims {
	if s == nil {
		return &jwt.IDTokenClaims{}
	}

	// Création de la structure claims
	claims := &jwt.IDTokenClaims{}

	// Subject (sub)
	claims.Subject = s.Subject

	// auth_time
	claims.AuthTime = s.AuthTime

	// acr
	if s.ACR != "" {
		claims.AuthenticationContextClassReference = s.ACR
	}

	// amr (stocké en JSON dans s.AMR)
	if s.AMR != nil {
		var amr []string
		_ = json.Unmarshal(s.AMR, &amr)
		claims.AuthenticationMethodsReferences = amr
	}

	// Extra claims (merge depuis s.Extra)
	claims.Extra = s.GetExtraClaims()

	// iat : temps d'émission = maintenant (ou si présent dans extra, on l'utilise)
	now := time.Now().UTC()

	claims.RequestedAt = now

	//algorithme de hash
	alg := "RS256"
	if s.Client.ID != uuid.Nil {
		alg = s.Client.GetRequestObjectSigningAlgorithm()
		if alg == "" {
			alg = "RS256"
		}
	}

	claims.JTI = uuid.New().String()

	return claims
}

// IDTokenHeaders construit et retourne les headers JWT (kid, typ, alg, etc).
// on essaye d'extraire un kid depuis le client (première JWK) sinon on utilise s.Extra["kid"].
func (s *Session) IDTokenHeaders() *jwt.Headers {
	if s == nil {
		return &jwt.Headers{}
	}

	headers := &jwt.Headers{}

	alg := "RS256"
	if s.Client.ID != uuid.Nil {
		a := s.Client.GetRequestObjectSigningAlgorithm()
		if a != "" {
			alg = a
		}
	}

	// kid : on essaye d'obtenir la première clé JWK du client
	kid := ""
	if s.Client.ID != uuid.Nil {
		jwks := s.Client.GetJSONWebKeys()
		if jwks != nil && len(jwks.Keys) > 0 {
			if jwks.Keys[0].KeyID != "" {
				kid = jwks.Keys[0].KeyID
			}
		}
	}

	extraMap := map[string]interface{}{
		"alg": alg,
		"typ": "JWT",
	}
	if kid != "" {
		extraMap["kid"] = kid
	}

	headers.Extra = extraMap
	return headers
}

func (s *Session) GetJWTHeader() *jwt.Headers {
	return s.IDTokenHeaders()
}
