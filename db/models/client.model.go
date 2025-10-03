package models

import (
	"encoding/json"
	"time"

	"github.com/dylEasydev/go-oauth2-easyclass/validators"
	"github.com/go-jose/go-jose/v3"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/ory/fosite"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// structure du client
type Client struct {
	ID     uuid.UUID `gorm:"primarykey;type:uuid;default:uuid_generate_v4()"`
	Active *bool     `gorm:"default:true"`

	//clés secret du client
	Secret string

	//listes des clés secrets de rotation
	RotatedSecrets pq.StringArray `gorm:"type:text[]"`

	Public *bool `gorm:"default:false"`

	//url de redirections
	RedirectURIs pq.StringArray `gorm:"type:text[]" validate:"required,urlallowed"`

	//Permissions demandé et accordé
	Scopes   pq.StringArray `gorm:"type:text[]"`
	Audience pq.StringArray `gorm:"type:text[]"`

	//grant de l'utilisateurs
	Grants        pq.StringArray `validate:"required,grantallowed"`
	ResponseTypes pq.StringArray `validate:"required,responseallowed"`

	//uri de ressources du client
	RequestURIs pq.StringArray `gorm:"type:text[]"`

	//modes de response "query" , "fragment" , "from_post"
	ResponseModes pq.StringArray `gorm:"type:text[]"`

	//methode d'authentification "client_secret_basic", "client_secret_post", "none", "private_key_jwt"
	TokenEndpointAuthMethod string `validate:"required,authmethodallowed"`

	RequestObjectSigningAlg           string `gorm:"type:text;default:'RS256'"`
	TokenEndpointAuthSigningAlgorithm string `gorm:"type:text;default:'RS256'"`

	//timestamps
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	//realtions avec information du client
	InfoClientID uuid.UUID  `gorm:"type:uuid;not null"`
	InfoClient   InfoClient `gorm:"foreignKey:InfoClientID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	//ensembel de clé public du client
	Keys []ClientKey `gorm:"foreignKey:ClientID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (Client) TableName() string {
	return "clients"
}

func (client *Client) BeforeSave(tx *gorm.DB) (err error) {
	// Validation
	if err = validators.ValidateStruct(client); err != nil {
		return
	}
	// Ne pas hasher si le client est public
	if client.Public != nil && *client.Public {
		client.Secret = ""
		client.RotatedSecrets = []string{}
		return nil
	}

	// Hasher le secret uniquement si ce n'est pas déjà hashé
	if _, err := bcrypt.Cost([]byte(client.Secret)); err != nil {
		hash, err := bcrypt.GenerateFromPassword([]byte(client.Secret), 10)
		if err != nil {
			return err
		}
		client.Secret = string(hash)
	}

	// Hasher les secrets de rotation uniquement si non hashés
	for i := range client.RotatedSecrets {
		if _, err := bcrypt.Cost([]byte(client.RotatedSecrets[i])); err != nil {
			hash, err := bcrypt.GenerateFromPassword([]byte(client.RotatedSecrets[i]), 10)
			if err != nil {
				return err
			}
			client.RotatedSecrets[i] = string(hash)
		}
	}

	return nil
}

func (client Client) GetID() string {
	return client.ID.String()
}

func (c *Client) GetHashedSecret() []byte {
	return []byte(c.Secret)
}

func (c *Client) GetRotatedHashes() [][]byte {
	var secrets [][]byte

	for _, secret := range c.RotatedSecrets {
		secrets = append(secrets, []byte(secret))
	}

	return secrets
}

func (c *Client) GetRedirectURIs() []string {
	var URIs []string

	for _, st := range c.RedirectURIs {
		URIs = append(URIs, st)
	}

	return URIs
}

func (c *Client) GetGrantTypes() fosite.Arguments {
	var Grants []string

	for _, st := range c.Grants {
		Grants = append(Grants, st)
	}

	return Grants
}

func (c *Client) GetResponseTypes() fosite.Arguments {
	var responses []string

	for _, st := range c.ResponseTypes {
		responses = append(responses, st)
	}

	return responses
}
func (c *Client) GetScopes() fosite.Arguments {
	var Scopes []string

	for _, st := range c.Scopes {
		Scopes = append(Scopes, st)
	}

	return Scopes
}
func (c *Client) IsPublic() bool {
	return c.Public != nil && *c.Public
}
func (c *Client) GetAudience() fosite.Arguments {
	var Audience []string

	for _, st := range c.Audience {
		Audience = append(Audience, st)
	}

	return Audience
}
func (c *Client) GetTokenEndpointAuthMethod() string {
	return c.TokenEndpointAuthMethod
}

func (c *Client) GetRequestURIs() []string {
	var requestURI []string

	for _, uri := range c.RequestURIs {
		requestURI = append(requestURI, uri)
	}

	return requestURI

}

func (c *Client) GetJSONWebKeys() *jose.JSONWebKeySet {
	keys := []jose.JSONWebKey{}
	for _, ck := range c.Keys {
		var jwk jose.JSONWebKey
		if err := json.Unmarshal(ck.JWK, &jwk); err == nil {
			keys = append(keys, jwk)
		}
	}
	return &jose.JSONWebKeySet{Keys: keys}
}

// URI vers le end-poin qui donne les clé JWK du client
func (c *Client) GetJSONWebKeysURI() string {
	return "/client/jwks/" + c.ID.String()
}

// Algorithme de signature pour Request Objects
func (c *Client) GetRequestObjectSigningAlgorithm() string {
	return c.RequestObjectSigningAlg
}

// Algorithme utilisé pour client_assertion (private_key_jwt)
func (c *Client) GetTokenEndpointAuthSigningAlgorithm() string {
	return c.TokenEndpointAuthSigningAlgorithm
}
