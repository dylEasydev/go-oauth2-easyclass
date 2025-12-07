package controller

import (
	"net/http"

	"github.com/dylEasydev/go-oauth2-easyclass/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-jose/go-jose/v3"
)

func JWKHandler(c *gin.Context) {
	publickey, err := utils.LoadPublicKey("public.key")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "erreur de lecture de la clé"})
	}

	jwKey := jose.JSONWebKey{
		Key:       publickey,
		KeyID:     "easy-class",
		Algorithm: "RS256",
		Use:       "sig",
	}
	c.JSON(http.StatusOK, jose.JSONWebKeySet{Keys: []jose.JSONWebKey{jwKey}})
}
