package middleware

import (
	"net/http"

	"github.com/dylEasydev/go-oauth2-easyclass/utils"
	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		if len(ctx.Errors) > 0 {
			err := ctx.Errors.Last().Err

			httpErr, ok := err.(utils.HttpErrorsInterface)
			if ok {
				ctx.JSON(httpErr.GetStatus(), gin.H{
					"sucess":  false,
					"message": httpErr.Error(),
				})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"sucess":  false,
				"message": err.Error(),
			})
		}
	}
}
