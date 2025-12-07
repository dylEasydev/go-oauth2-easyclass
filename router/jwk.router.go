package router

import (
	"github.com/dylEasydev/go-oauth2-easyclass/controller"
)

func (r *Router) JWKRouter() {
	jwkGroup := r.Server.Group("/key")

	{
		jwkGroup.GET("/jws.json", controller.JWKHandler)
	}
}
