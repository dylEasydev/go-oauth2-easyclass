package router

import (
	"github.com/dylEasydev/go-oauth2-easyclass/middleware"
	"github.com/dylEasydev/go-oauth2-easyclass/utils"
)

func (r *router) SignRouter() {
	signGroup := r.Server.Group("/sign")
	publicKey, err := utils.LoadPublicKey("public")
	if err != nil {
		panic("impossible de lire la clé public")
	}

	{
		signGroup.POST("/teacher", r.StoreRequest.SignTeacher)
		signGroup.POST("/student", r.StoreRequest.SignStudent)
		signGroup.POST("/admin", middleware.AuthMiddleware(publicKey), r.StoreRequest.SignStudent)
	}
}
