package router

import (
	"github.com/dylEasydev/go-oauth2-easyclass/controller"
	"github.com/dylEasydev/go-oauth2-easyclass/db"
	"github.com/gin-gonic/gin"
)

type Router struct {
	Server *gin.Engine
	Store  *db.Store
}

func (r *Router) IndexRouter() {
	indexGroup := r.Server.Group("/")
	{
		indexGroup.GET("/", controller.IndexHanler)
	}
}
