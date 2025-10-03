package utils

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"strings"
)

// oidcHash calcule le at_hash ou c_hash pour OIDC :
//   - calcule le hash correspondant (SHA256/384/512 depending on alg)
//   - prend la moitié gauche des octets
//   - base64url encode (sans padding)
//
// Renvoie la string base64url ou une erreur.
func OidcHash(value string, alg string) (string, error) {
	var hashed []byte

	// choisir la fonction de hash en fonction de l'algorithme (RS256 -> sha256, RS384->sha384, RS512->sha512)
	alg = strings.ToUpper(alg)
	switch {
	case strings.HasSuffix(alg, "256"):
		h := sha256.Sum256([]byte(value))
		hashed = h[:]
	case strings.HasSuffix(alg, "384"):
		h := sha512.Sum384([]byte(value))
		hashed = h[:]
	case strings.HasSuffix(alg, "512"):
		h := sha512.Sum512([]byte(value))
		hashed = h[:]
	default:
		// par défaut sha256
		h := sha256.Sum256([]byte(value))
		hashed = h[:]
	}

	// prendre la moitié gauche des octets
	half := len(hashed) / 2
	left := hashed[:half]

	// base64url sans padding
	encoded := base64.RawURLEncoding.EncodeToString(left)
	return encoded, nil
}
