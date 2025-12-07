package router

import (
	"github.com/dylEasydev/go-oauth2-easyclass/controller"
	"github.com/dylEasydev/go-oauth2-easyclass/db"
	"github.com/gin-gonic/gin"
)

type router struct {
	Server       *gin.Engine
	Store        *db.Store
	StoreRequest *controller.StoreRequest
}

func NewRouter(server *gin.Engine, store *db.Store) *router {
	return &router{
		Server: server,
		Store:  store,
		StoreRequest: &controller.StoreRequest{
			Store: store,
		},
	}
}

func (r *router) IndexRouter() {
	indexGroup := r.Server.Group("/")
	{
		indexGroup.GET("/", controller.IndexHanler)
	}
}
