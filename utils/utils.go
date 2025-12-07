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

func LoadPrivateKey(fileName string) (*rsa.PrivateKey, error) {
	baseDir, _ := os.Getwd()
	fullPath := path.Join(baseDir, "key/", fileName+".json")
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("erreur ouverture fichier: %w", err)
	}
	block, _ := pem.Decode(data)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("échec du décodage de la clé privé")
	}
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

func LoadPublicKey(fileName string) (*rsa.PublicKey, error) {
	baseDir, _ := os.Getwd()
	fullPath := path.Join(baseDir, "key/", fileName+".json")
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
