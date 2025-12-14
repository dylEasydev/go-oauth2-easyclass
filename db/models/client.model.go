package models

//packages models

import (
	"time"

	"github.com/dylEasydev/go-oauth2-easyclass/utils"
	"github.com/dylEasydev/go-oauth2-easyclass/validators"
	"github.com/go-jose/go-jose/v3"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/ory/fosite"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// db models Client pour OIDC
type Client struct {
	ID     uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Active *bool     `gorm:"default:true"`

	//clés secret du client
	Secret string `gorm:"not null"`

	//listes des clés secrets de rotation
	RotatedSecrets pq.StringArray `gorm:"type:text[]"`

	//client public ou privé
	Public *bool `gorm:"default:false"`

	//url de redirections
	RedirectURIs pq.StringArray `gorm:"type:text[]" validate:"required,urlallowed"`

	//Permissions demandé et accordé
	Scopes   pq.StringArray `gorm:"type:text[]"`
	Audience pq.StringArray `gorm:"type:text[]"`

	//grant du client
	Grants pq.StringArray `gorm:"type:text[]" validate:"required,grantallowed"`

	//types de la response
	ResponseTypes pq.StringArray `gorm:"type:text[]" validate:"required,responseallowed"`

	//uri de ressources du client
	RequestURIs pq.StringArray `gorm:"type:text[]"`

	//modes de response "query" , "fragment" , "from_post"
	ResponseModes pq.StringArray `gorm:"type:text[]"`

	//methode d'authentification "client_secret_basic", "client_secret_post", "none", "private_key_jwt"
	TokenEndpointAuthMethod string `validate:"required,authmethodallowed"`

	// algorithme de signature des jetons assertion
	RequestObjectSigningAlg           string `gorm:"type:text;default:'RS256'"`
	TokenEndpointAuthSigningAlgorithm string `gorm:"type:text;default:'RS256'"`

	//timestamps
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	//realtions avec information du client
	InfoClientID uuid.UUID  `gorm:"type:uuid;not null"`
	InfoClient   InfoClient `gorm:"foreignKey:InfoClientID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	//ensemble des clé public du client
	Keys []ClientKey `gorm:"foreignKey:ClientID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// implementation de interface Tabler(pour le nom de la table)
func (Client) TableName() string {
	return "clients"
}

// hooks avant la sauvegarde du client
func (client *Client) BeforeSave(tx *gorm.DB) error {
	// Validation
	if err := validators.ValidateStruct(client); err != nil {
		return err
	}
	// Ne pas hasher si le client est public
	if client.Public != nil && *client.Public {
		client.Secret = ""
		client.RotatedSecrets = []string{}
		return nil
	}

	// Hashe le secret uniquement si ce n'est pas déjà hashé
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

// implementation de interface client de Fosite

// récupère id du client
func (client Client) GetID() string {
	return client.ID.String()
}

// récupère le secret du client
func (c *Client) GetHashedSecret() []byte {
	return []byte(c.Secret)
}

// récupères les secrets de rotation du client
func (c *Client) GetRotatedHashes() [][]byte {
	var secrets [][]byte

	for _, secret := range c.RotatedSecrets {
		secrets = append(secrets, []byte(secret))
	}

	return secrets
}

func (c *Client) VerifySecret(plain string) bool {
	// Ne pas vérifier pour les clients publics
	if c.IsPublic() {
		return false
	}

	// Si la méthode est "none", pas de verification de secret
	if c.TokenEndpointAuthMethod == "none" {
		return false
	}

	// Vérifie le secret courant s'il existe
	if c.Secret != "" {
		if err := bcrypt.CompareHashAndPassword([]byte(c.Secret), []byte(plain)); err == nil {
			return true
		}
	}

	// Vérifie les secrets de rotation
	for _, s := range c.RotatedSecrets {
		if s == "" {
			continue
		}
		if err := bcrypt.CompareHashAndPassword([]byte(s), []byte(plain)); err == nil {
			return true
		}
	}

	return false
}

// récupère les url de redirection
func (c *Client) GetRedirectURIs() []string {
	var URIs []string

	for _, st := range c.RedirectURIs {
		URIs = append(URIs, st)
	}

	return URIs
}

// récupères les grant_type du client
func (c *Client) GetGrantTypes() fosite.Arguments {
	var Grants []string

	for _, st := range c.Grants {
		Grants = append(Grants, st)
	}

	return Grants
}

// récupères les responses types du client
func (c *Client) GetResponseTypes() fosite.Arguments {
	var responses []string

	for _, st := range c.ResponseTypes {
		responses = append(responses, st)
	}

	return responses
}

// récupères les scope(permissions) du client
func (c *Client) GetScopes() fosite.Arguments {
	var Scopes []string

	for _, st := range c.Scopes {
		Scopes = append(Scopes, st)
	}

	return Scopes
}

// verifie si un client est  public ou non
func (c *Client) IsPublic() bool {
	return c.Public != nil && *c.Public
}

// récupères les permissions accorder au client
func (c *Client) GetAudience() fosite.Arguments {
	var Audience []string

	for _, st := range c.Audience {
		Audience = append(Audience, st)
	}

	return Audience
}

// récupère la méthodes authentifiaction du client
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

// récupère l'ensemble des clés public du client
func (c *Client) GetJSONWebKeys() *jose.JSONWebKeySet {
	keys := []jose.JSONWebKey{}
	for _, ck := range c.Keys {
		keys = append(keys, jose.JSONWebKey(ck.JWK))
	}
	return &jose.JSONWebKeySet{Keys: keys}
}

// URI vers le end-point qui donne les clé JWK du client
func (c *Client) GetJSONWebKeysURI() string {
	return utils.URL_Host + "/keys/client/jwks/" + c.ID.String()
}

// Algorithme de signature pour Request Objects
func (c *Client) GetRequestObjectSigningAlgorithm() string {
	return c.RequestObjectSigningAlg
}

// Algorithme utilisé pour client_assertion (private_key_jwt)
func (c *Client) GetTokenEndpointAuthSigningAlgorithm() string {
	return c.TokenEndpointAuthSigningAlgorithm
}
