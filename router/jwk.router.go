package router

import (
	"github.com/dylEasydev/go-oauth2-easyclass/controller"
)

func (r *router) JWKRouter() {
	jwkGroup := r.Server.Group("/keys")

	{
		jwkGroup.GET("/clients/jwks/:id", r.StoreRequest.ClientJWKHanler)
		jwkGroup.GET("/jwks.json", controller.JWKHandler)
	}
}
