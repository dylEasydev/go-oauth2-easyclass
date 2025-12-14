package utils

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path"
)

func PtrBool(val bool) *bool {
	return &val
}

// intersection des scopes de client et d'utilisateur
func IntersectScopes(clientScopes, userScopes []string) []string {
	result := make([]string, len(clientScopes))
	setScopes := make(map[string]bool, len(userScopes))

	for _, scope := range userScopes {
		setScopes[scope] = true
	}

	for _, scope := range clientScopes {
		if scope == "openid" || setScopes[scope] {
			result = append(result, scope)
		}
	}
	return result
}

// lecture de clé privé RSA
func LoadPrivateKey(fileName string) (*rsa.PrivateKey, error) {
	baseDir, _ := os.Getwd()
	fullPath := path.Join(baseDir, "key/", fileName+".key")
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("erreur ouverture fichier: %w", err)
	}
	block, _ := pem.Decode(data)
	if block == nil || block.Type != "PRIVATE KEY" {
		return nil, fmt.Errorf("échec du décodage de la clé privé")
	}
	priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	switch priv := priv.(type) {
	case *rsa.PrivateKey:
		return priv, nil
	default:
		return nil, fmt.Errorf("la clé publicn'est pas de type RSA ")
	}
}

// lecture de clé public RSA
func LoadPublicKey(fileName string) (*rsa.PublicKey, error) {
	baseDir, _ := os.Getwd()
	fullPath := path.Join(baseDir, "key/", fileName+".key")
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("erreur ouverture fichier: %w", err)
	}
	block, _ := pem.Decode(data)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, fmt.Errorf("échec du décodage de la clé public")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	switch pub := pub.(type) {
	case *rsa.PublicKey:
		return pub, nil
	default:
		return nil, fmt.Errorf("la clé publicn'est pas de type RSA ")
	}
}
