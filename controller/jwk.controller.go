package controller

import (
	"errors"
	"net/http"

	"github.com/dylEasydev/go-oauth2-easyclass/db"
	"github.com/dylEasydev/go-oauth2-easyclass/db/models"
	"github.com/dylEasydev/go-oauth2-easyclass/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-jose/go-jose/v3"
	"github.com/ory/fosite"
)

type StoreRequest struct {
	Store *db.Store
}

type IDUri struct {
	ID string `uri:"id" binding:"required,uuid"`
}

func (s *StoreRequest) ClientJWKHanler(c *gin.Context) {
	var idClient IDUri
	context := c.Request.Context()

	if err := c.ShouldBindUri(&idClient); err != nil {
		httpErr := utils.HttpErrors{Status: http.StatusBadRequest, Message: err.Error()}
		c.Error(&httpErr)
		return
	}

	client, err := s.Store.GetClient(context, idClient.ID)
	if err != nil {
		if errors.Is(err, fosite.ErrNotFound) {
			httpErr := utils.HttpErrors{Status: http.StatusNotFound, Message: err.Error()}
			c.Error(&httpErr)
			return
		}
		c.Error(err)
		return
	}

	clientJWK := client.(*models.Client)

	c.JSON(http.StatusOK, clientJWK.GetJSONWebKeys())
}

func JWKHandler(c *gin.Context) {
	publickey, err := utils.LoadPublicKey("public.key")
	if err != nil {
		httpErr := utils.HttpErrors{Status: http.StatusInternalServerError, Message: err.Error()}
		c.Error(&httpErr)
		return
	}

	jwKey := jose.JSONWebKey{
		Key:       publickey,
		KeyID:     "easy-class",
		Algorithm: "RS256",
		Use:       "sig",
	}
	c.JSON(http.StatusOK, jose.JSONWebKeySet{Keys: []jose.JSONWebKey{jwKey}})
}
