package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func IndexHanler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"messages": "Bienvenu sur oauth2 easyclass",
	})
}
