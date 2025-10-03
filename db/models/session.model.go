package models

import (
	"encoding/json"
	"time"

	"github.com/dylEasydev/go-oauth2-easyclass/utils"
	"github.com/google/uuid"
	"github.com/mohae/deepcopy"
	"github.com/ory/fosite"
	"github.com/ory/fosite/token/jwt"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// structure de session basée sur fosite
type Session struct {
	ID uuid.UUID `gorm:"primarykey;type:uuid;default:uuid_generate_v4()"`

	Username  string
	Subject   string
	ExpiresAt datatypes.JSON `gorm:"type:jsonb"`

	//OIDC specification
	Nonce    string
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

// GetUsername returns the username, if set. This is optional and only used during token introspection.
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
		_ = json.Unmarshal(s.Extra, &extra)
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

func (s *Session) IDTokenClaims() *jwt.IDTokenClaims {
	if s == nil {
		return &jwt.IDTokenClaims{}
	}

	// Création de la structure claims
	claims := &jwt.IDTokenClaims{}

	// Issuer : on cherche dans s.Extra["iss"] sinon vide (à configurer)
	if extra := s.GetExtraClaims(); extra != nil {
		if iss, ok := extra["iss"].(string); ok && iss != "" {
			claims.Issuer = iss
		}
	}

	// Subject (sub)
	claims.Subject = s.Subject

	// Audience (aud) : on préfère l'audience du client si présente,
	// sinon on met l'ID du client comme audience.
	if s.ClientID != uuid.Nil {
		aud := s.Client.GetAudience()
		if len(aud) > 0 {
			for _, a := range aud {
				claims.Audience = append(claims.Audience, a)
			}
		} else {
			claims.Audience = []string{s.ClientID.String()}
		}
	}

	// Nonce
	claims.Nonce = s.Nonce

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
	if claims.IssuedAt.IsZero() {
		claims.IssuedAt = now
	}

	// requested at (rat) : si fourni dans extra on l'utilise, sinon now
	if rat, ok := claims.Extra["rat"].(time.Time); ok {
		claims.RequestedAt = rat
	} else if str, ok := claims.Extra["rat"].(string); ok {
		// tentative de parsing si l'utilisateur a stocké une string
		if t, err := time.Parse(time.RFC3339, str); err == nil {
			claims.RequestedAt = t
		} else {
			claims.RequestedAt = now
		}
	} else {
		claims.RequestedAt = now
	}

	// exp : on récupère l'expiration depuis la session (si elle existe)
	if exp := s.GetExpiresAt(fosite.IDToken); !exp.IsZero() {
		claims.ExpiresAt = exp
	} else if exp := s.GetExpiresAt(fosite.AccessToken); !exp.IsZero() {
		claims.ExpiresAt = exp
	} else if exp := s.GetExpiresAt(fosite.AuthorizeCode); !exp.IsZero() {
		claims.ExpiresAt = exp
	}

	// at_hash / c_hash : si access_token ou code sont fournis dans Extra,
	// on calcule l'hash conformément à l'algorithme de signature (RS256/ES256/etc).
	// alg preferé : essayer extra["alg"], sinon default "RS256".
	alg := "RS256"
	if a, ok := claims.Extra["alg"].(string); ok && a != "" {
		alg = a
	} else {
		// si le client a une préférence d'alg (ex: RequestObjectSigningAlg), on peut l'utiliser
		if s.Client.ID != uuid.Nil {
			alg = s.Client.GetRequestObjectSigningAlgorithm()
			if alg == "" {
				alg = "RS256"
			}
		}
	}

	// calcul de c_hash si code dans extra
	if code, ok := claims.Extra["code"].(string); ok && code != "" {
		if h, err := utils.OidcHash(code, alg); err == nil {
			claims.CodeHash = h
		}
	}

	// jti : si fourni dans extra ou générer localement
	if jti, ok := claims.Extra["jti"].(string); ok && jti != "" {
		claims.JTI = jti
	}

	return claims
}

// IDTokenHeaders construit et retourne les headers JWT (kid, typ, alg, etc).
// on essaye d'extraire un kid depuis le client (première JWK) sinon on utilise s.Extra["kid"].
func (s *Session) IDTokenHeaders() *jwt.Headers {
	if s == nil {
		return &jwt.Headers{}
	}

	headers := &jwt.Headers{}

	// alg : on tente de récupérer depuis extra ou depuis la config du client
	alg := "RS256"
	if extra := s.GetExtraClaims(); extra != nil {
		if a, ok := extra["alg"].(string); ok && a != "" {
			alg = a
		}
	}

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
	// fallback : extra["kid"]
	if kid == "" {
		if extra := s.GetExtraClaims(); extra != nil {
			if k, ok := extra["kid"].(string); ok && k != "" {
				kid = k
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
