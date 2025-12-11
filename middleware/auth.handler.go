package middleware

import (
	"crypto/rsa"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/cristalhq/jwt/v4"
	"github.com/gin-gonic/gin"
	fosite_jwt "github.com/ory/fosite/token/jwt"
)

func AuthMiddleware(publicKey *rsa.PublicKey) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": "entête d'authorization non fournis",
				"success": false,
			})
			ctx.Abort()
			return
		}
		partsToken := strings.Split(authHeader, " ")
		if len(partsToken) != 2 || partsToken[0] != "Bearer" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": "pas de Bearer Token fournis ",
				"success": false,
			})
			ctx.Abort()
			return
		}

		tokenString := partsToken[1]

		verifier, err := jwt.NewVerifierRS("RS256", publicKey)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "erreur au niveau du serveur ",
				"success": false,
			})
			ctx.Abort()
			return
		}
		var claims = fosite_jwt.JWTClaims{}
		token, err := jwt.Parse([]byte(tokenString), verifier)

		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": "mauvais jeton fournis ",
				"success": false,
			})
			ctx.Abort()
			return
		}
		err = json.Unmarshal(token.Claims(), &claims)
		if err != nil {

		}
		ctx.Set("claims", claims)
		ctx.Next()
	}
}
