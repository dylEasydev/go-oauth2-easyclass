package router

import (
	"github.com/dylEasydev/go-oauth2-easyclass/controller"
	"github.com/dylEasydev/go-oauth2-easyclass/provider"
	"github.com/dylEasydev/go-oauth2-easyclass/utils"
)

func (r *router) OIDCRouter() {
	privateKey, err := utils.LoadPrivateKey("private.key")
	if err != nil {
		panic("impossible de lire les clé de signature")
	}
	provider := provider.InitProvider(r.Store, privateKey)
	auth := controller.NewAuth(provider, r.Store)
	oidcGroup := r.Server.Group("/oidc")

	{
		oidcGroup.POST("/authorize", auth.AuthorizeHandler)
		oidcGroup.GET("/authrize", auth.AuthorizeHandler)
		oidcGroup.POST("/token", auth.TokenHandler)
		oidcGroup.POST("/revoke", auth.RevokeHandler)
		oidcGroup.POST("/par", auth.PARRequestHandler)
		oidcGroup.POST("/introspect", auth.IntrospectionHandler)
	}
}
