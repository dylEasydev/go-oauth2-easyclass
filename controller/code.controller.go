package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/dylEasydev/go-oauth2-easyclass/db/interfaces"
	"github.com/dylEasydev/go-oauth2-easyclass/db/service"
	"github.com/dylEasydev/go-oauth2-easyclass/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
)

type CodeBody struct {
	Code string `form:"codeverif" json:"codeverif" binding:"required,min=6"`
}

type CodeUri struct {
	UserName  string `uri:"name" binding:"required,name"`
	TableName string `uri:"table" binding:"required,tableName"`
}

func (s *StoreRequest) VerifCode(ctx *gin.Context) {
	var idUser IDUri
	var codeBody CodeBody

	context := ctx.Request.Context()

	if err := ctx.ShouldBindUri(&idUser); err != nil {
		httpErr := utils.HttpErrors{Status: http.StatusBadRequest, Message: err.Error()}
		ctx.Error(&httpErr)
		return
	}
	id, err := uuid.Parse(idUser.ID)
	if err != nil {
		httpErr := utils.HttpErrors{Status: http.StatusBadRequest, Message: err.Error()}
		ctx.Error(&httpErr)
		return
	}

	if err := ctx.ShouldBindWith(&codeBody, binding.JSON); err != nil {
		httpErr := utils.HttpErrors{Status: http.StatusBadRequest, Message: err.Error()}
		ctx.Error(&httpErr)
		return
	}
	codeHash := utils.GenerateHash(codeBody.Code)

	codeservice := service.InitCodeService(context, s.Store.GetDb())

	codeVerif, err := codeservice.FindCode(codeHash, id)
	if err != nil {
		if errors.Is(err, service.ErrNotCode) {
			httpErr := utils.HttpErrors{Status: http.StatusBadRequest, Message: err.Error()}
			ctx.Error(&httpErr)
			return
		}
		httpErr := utils.HttpErrors{Status: http.StatusInternalServerError, Message: err.Error()}
		ctx.Error(&httpErr)
		return
	}

	if codeVerif.IsUsed() || codeVerif.IsExpired() {
		httpErr := utils.HttpErrors{Status: http.StatusUnauthorized, Message: "code de vérification non valide"}
		ctx.Error(&httpErr)
		return
	}

	user, err := codeVerif.GetForeign(s.Store.GetDb())
	if err != nil {
		httpErr := utils.HttpErrors{Status: http.StatusInternalServerError, Message: err.Error()}
		ctx.Error(&httpErr)
		return
	}

	userTemp, ok := user.(interfaces.UserTempInterafce)
	if !ok {
		httpErr := utils.HttpErrors{Status: http.StatusBadRequest, Message: "end-point réservé au utilisateur en attente"}
		ctx.Error(&httpErr)
		return
	}

	userTempservice := service.InitUserTempService(s.Store.GetDb())
	if err := userTempservice.SaveUser(userTemp); err != nil {
		if !errors.Is(err, service.ErrDestroy) {
			httpErr := utils.HttpErrors{Status: http.StatusInternalServerError, Message: err.Error()}
			ctx.Error(&httpErr)
			return
		}
	}

	if err := codeVerif.MarkUsed(s.Store.GetDb()); err != nil {
		httpErr := utils.HttpErrors{Status: http.StatusInternalServerError, Message: err.Error()}
		ctx.Error(&httpErr)
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"sucess":  true,
		"message": fmt.Sprintf("Bienvenu utilisateur @%s", userTemp.GetName()),
		"data":    userTemp,
	})
}

func (s *StoreRequest) RestartCode(ctx *gin.Context) {
	var codeUri CodeUri

	context := ctx.Request.Context()

	if err := ctx.ShouldBindUri(&codeUri); err != nil {
		httpErr := utils.HttpErrors{Status: http.StatusBadRequest, Message: err.Error()}
		ctx.Error(&httpErr)
		return
	}

	codeservice := service.InitCodeService(context, s.Store.GetDb())
	user, err := codeservice.GetForeignByName(codeUri.TableName, codeUri.UserName)
	if err != nil {
		httpErr := utils.HttpErrors{Status: http.StatusInternalServerError, Message: err.Error()}
		ctx.Error(&httpErr)
		return
	}
	code, err := codeservice.FindCodeTable(codeUri.TableName, user.GetId())
	if err != nil {
		if errors.Is(err, service.ErrNotCode) {
			httpErr := utils.HttpErrors{Status: http.StatusBadRequest, Message: err.Error()}
			ctx.Error(&httpErr)
			return
		}
		httpErr := utils.HttpErrors{Status: http.StatusInternalServerError, Message: err.Error()}
		ctx.Error(&httpErr)
		return
	}
	if err := codeservice.UpdateCodeVerif(user, code.Code); err != nil {
		httpErr := utils.HttpErrors{Status: http.StatusInternalServerError, Message: err.Error()}
		ctx.Error(&httpErr)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"sucess":  true,
		"message": fmt.Sprintf("verifier votre mail %s @%s", user.GetMail(), user.GetName()),
	})
}
