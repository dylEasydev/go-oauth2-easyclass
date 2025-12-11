package controller

import (
	"errors"
	"fmt"
	"net/http"
	"slices"

	"github.com/dylEasydev/go-oauth2-easyclass/db/models"
	"github.com/dylEasydev/go-oauth2-easyclass/db/service"
	"github.com/dylEasydev/go-oauth2-easyclass/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/ory/fosite/token/jwt"
	"gorm.io/gorm"
)

type TeacherBody struct {
	UserName    string `form:"name" json:"name" binding:"required,name"`
	Password    string `form:"password" json:"password" binding:"required,min=8,password"`
	Email       string `form:"email" json:"email" binding:"required,email"`
	SubjectName string `form:"subject" json:"subject" binding:"required,name"`
}

type UserBody struct {
	UserName string `form:"name" json:"name" binding:"required,name"`
	Password string `form:"password" json:"password" binding:"required,min=8,password"`
	Email    string `form:"email" json:"email" binding:"required,email"`
}

func (s *StoreRequest) SignTeacher(ctx *gin.Context) {
	var bodyTeacher TeacherBody

	context := ctx.Request.Context()

	if err := ctx.ShouldBindWith(&bodyTeacher, binding.Query); err != nil {
		httpErr := utils.HttpErrors{Status: http.StatusBadRequest, Message: err.Error()}
		ctx.Error(&httpErr)
		return
	}

	//recherche parmis la table des utilisateur
	userFind, err := service.FindUserByName[models.User](context, s.Store.GetDb(), bodyTeacher.UserName, bodyTeacher.Email)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			httpErr := utils.HttpErrors{Message: err.Error(), Status: http.StatusInternalServerError}
			ctx.Error(&httpErr)
			return
		}
	}
	if userFind != nil {
		httpErr := utils.HttpErrors{Message: "utilisateurs possède déjà un compte ", Status: http.StatusBadRequest}
		ctx.Error(&httpErr)
		return
	}

	// recherche parmi la table des enseignant  temportaire

	teacherFind, err := service.FindUserByName[models.TeacherWaiting](context, s.Store.GetDb(), bodyTeacher.UserName, bodyTeacher.Email)

	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			httpErr := utils.HttpErrors{Message: err.Error(), Status: http.StatusInternalServerError}
			ctx.Error(&httpErr)
			return
		}
	}

	if teacherFind != nil {
		httpErr := utils.HttpErrors{Message: "utilisateurs possède déjà un compte ", Status: http.StatusBadRequest}
		ctx.Error(&httpErr)
		return
	}
	teacherService := service.TeacherService{Db: s.Store.GetDb(), Ctx: &context}
	data := service.TeacherBody{
		UserBody: service.UserBody{
			Name:     bodyTeacher.UserName,
			Email:    bodyTeacher.Email,
			Password: bodyTeacher.Password,
		},
		Subject: bodyTeacher.SubjectName,
	}
	newTeacher, err := teacherService.CreateUser(&data)
	if err != nil {
		httpErr := utils.HttpErrors{Message: err.Error(), Status: http.StatusInternalServerError}
		ctx.Error(&httpErr)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"sucess":  true,
		"message": fmt.Sprintf("verifier votre mail %s @%s", newTeacher.Email, newTeacher.UserName),
	})
}

func (s *StoreRequest) SignStudent(ctx *gin.Context) {
	var bodyStudent UserBody

	context := ctx.Request.Context()

	if err := ctx.ShouldBindWith(&bodyStudent, binding.Query); err != nil {
		httpErr := utils.HttpErrors{Status: http.StatusBadRequest, Message: err.Error()}
		ctx.Error(&httpErr)
		return
	}

	//recherche parmis la table des utilisateur
	userFind, err := service.FindUserByName[models.User](context, s.Store.GetDb(), bodyStudent.UserName, bodyStudent.Email)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			httpErr := utils.HttpErrors{Message: err.Error(), Status: http.StatusInternalServerError}
			ctx.Error(&httpErr)
			return
		}
	}
	if userFind != nil {
		httpErr := utils.HttpErrors{Message: "utilisateurs possède déjà un compte ", Status: http.StatusBadRequest}
		ctx.Error(&httpErr)
		return
	}

	// recherche parmi la table des enseignant  temportaire

	teacherFind, err := service.FindUserByName[models.TeacherWaiting](context, s.Store.GetDb(), bodyStudent.UserName, bodyStudent.Email)

	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			httpErr := utils.HttpErrors{Message: err.Error(), Status: http.StatusInternalServerError}
			ctx.Error(&httpErr)
			return
		}
	}

	if teacherFind != nil {
		httpErr := utils.HttpErrors{Message: "utilisateurs possède déjà un compte ", Status: http.StatusBadRequest}
		ctx.Error(&httpErr)
		return
	}
	studentService := service.StudentService{Db: s.Store.GetDb(), Ctx: &context}
	data := service.UserBody{
		Name:     bodyStudent.UserName,
		Email:    bodyStudent.Email,
		Password: bodyStudent.Password,
	}
	newStudent, err := studentService.CreateUser(&data)
	if err != nil {
		httpErr := utils.HttpErrors{Message: err.Error(), Status: http.StatusInternalServerError}
		ctx.Error(&httpErr)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"sucess":  true,
		"message": fmt.Sprintf("verifier votre mail %s @%s", newStudent.Email, newStudent.UserName),
	})
}

func (s *StoreRequest) SignAdmin(ctx *gin.Context) {
	claims, ok := ctx.Get("claims")
	if !ok {
		httpErr := utils.HttpErrors{Status: http.StatusUnauthorized, Message: "vous n'avez pas fourni de jeton JWT "}
		ctx.Error(&httpErr)
		return
	}
	convertClaims := claims.(jwt.JWTClaims)
	if !slices.Contains(convertClaims.Scope, "admin.created") && !slices.Contains(convertClaims.Scope, "admin.*") {
		httpErr := utils.HttpErrors{Status: http.StatusForbidden, Message: "vous n'avez pas les autorisation pour créer un administarteur "}
		ctx.Error(&httpErr)
		return
	}
	var bodyUser UserBody

	context := ctx.Request.Context()

	if err := ctx.ShouldBindWith(&bodyUser, binding.Query); err != nil {
		httpErr := utils.HttpErrors{Status: http.StatusBadRequest, Message: err.Error()}
		ctx.Error(&httpErr)
		return
	}

	//recherche parmis la table des utilisateur
	userFind, err := service.FindUserByName[models.User](context, s.Store.GetDb(), bodyUser.UserName, bodyUser.Email)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			httpErr := utils.HttpErrors{Message: err.Error(), Status: http.StatusInternalServerError}
			ctx.Error(&httpErr)
			return
		}
	}
	if userFind != nil {
		httpErr := utils.HttpErrors{Message: "utilisateurs possède déjà un compte ", Status: http.StatusBadRequest}
		ctx.Error(&httpErr)
		return
	}

	// recherche parmi la table des enseignant  temportaire

	teacherFind, err := service.FindUserByName[models.TeacherWaiting](context, s.Store.GetDb(), bodyUser.UserName, bodyUser.Email)

	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			httpErr := utils.HttpErrors{Message: err.Error(), Status: http.StatusInternalServerError}
			ctx.Error(&httpErr)
			return
		}
	}

	if teacherFind != nil {
		httpErr := utils.HttpErrors{Message: "utilisateurs possède déjà un compte ", Status: http.StatusBadRequest}
		ctx.Error(&httpErr)
		return
	}
	userService := service.UserService{Db: s.Store.GetDb(), Ctx: &context}
	data := service.UserBody{
		Name:     bodyUser.UserName,
		Email:    bodyUser.Email,
		Password: bodyUser.Password,
	}
	newUser, err := userService.CreateUser(&data)
	if err != nil {
		httpErr := utils.HttpErrors{Message: err.Error(), Status: http.StatusInternalServerError}
		ctx.Error(&httpErr)
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"sucess":  true,
		"message": fmt.Sprintf("Bienvenu admin  @%s", newUser.UserName),
		"data":    newUser,
	})
}
