package controller

import (
	"github.com/dylEasydev/go-oauth2-easyclass/db"
	"github.com/dylEasydev/go-oauth2-easyclass/db/models"
	"github.com/dylEasydev/go-oauth2-easyclass/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/ory/fosite"
)

type Auth struct {
	provider fosite.OAuth2Provider
	store    *db.Store
}

type Authorize struct {
	UserName string   `form:"username" json:"username" binding:"required,name"`
	Password string   `form:"password" json:"password" binding:"required,min=8,password"`
	Scopes   []string `form:"scopes" json:"scopes"`
}

func NewAuth(provider fosite.OAuth2Provider, store *db.Store) *Auth {
	return &Auth{
		provider: provider,
		store:    store,
	}
}

func (a *Auth) AuthorizeHandler(c *gin.Context) {
	ctx := c.Request.Context()

	authorizeRequest, err := a.provider.NewAuthorizeRequest(ctx, c.Request)
	if err != nil {
		a.provider.WriteAuthorizeError(ctx, c.Writer, authorizeRequest, err)
		return
	}

	form := Authorize{}
	if err := c.ShouldBindWith(&form, binding.Query); err != nil {
		a.provider.WriteAuthorizeError(ctx, c.Writer, authorizeRequest, fosite.ErrInvalidRequest.WithHint(err.Error()))
		return
	}

	err = a.store.Authenticate(ctx, form.UserName, form.Password)
	if err != nil {
		a.provider.WriteAuthorizeError(ctx, c.Writer, authorizeRequest, fosite.ErrAccessDenied.WithHint(err.Error()))
		return
	}

	user, err := a.store.GetUser(ctx, form.UserName)
	if err != nil {
		a.provider.WriteAuthorizeError(ctx, c.Writer, authorizeRequest, err)
		return
	}

	userScopes := make([]string, len(user.Role.Scopes))
	for _, scopes := range user.Role.Scopes {
		userScopes = append(userScopes, scopes.ScopeName)
	}

	grantScopes := utils.IntersectScopes(form.Scopes, userScopes)

	for _, scope := range grantScopes {
		authorizeRequest.GrantScope(scope)
	}

	extra := map[string]any{
		"scopes":  grantScopes,
		"user_id": user.ID,
	}
	session, err := models.NewSession(ctx, authorizeRequest.GetClient().GetID(), user.ID.String(), user.UserName, user.UserName, extra)

	if err != nil {
		a.provider.WriteAuthorizeError(ctx, c.Writer, authorizeRequest, err)
		return
	}
	response, err := a.provider.NewAuthorizeResponse(ctx, authorizeRequest, session)
	if err != nil {
		a.provider.WriteAuthorizeError(ctx, c.Writer, authorizeRequest, err)
		return
	}

	a.provider.WriteAuthorizeResponse(ctx, c.Writer, authorizeRequest, response)
}

func (a *Auth) TokenHandler(c *gin.Context) {
	ctx := c.Request.Context()

	accessRequest, err := a.provider.NewAccessRequest(ctx, c.Request, new(models.Session))
	if err != nil {
		a.provider.WriteAccessError(ctx, c.Writer, accessRequest, err)
		return
	}

	if accessRequest.GetGrantTypes().ExactOne("client_credentials") {
		for _, scope := range accessRequest.GetRequestedScopes() {
			accessRequest.GrantScope(scope)
		}
	}

	response, err := a.provider.NewAccessResponse(ctx, accessRequest)
	if err != nil {
		a.provider.WriteAccessError(ctx, c.Writer, accessRequest, err)
		return
	}

	a.provider.WriteAccessResponse(ctx, c.Writer, accessRequest, response)
}

func (a *Auth) PARRequestHandler(c *gin.Context) {
	ctx := c.Request.Context()

	parRequest, err := a.provider.NewPushedAuthorizeRequest(ctx, c.Request)
	if err != nil {
		a.provider.WritePushedAuthorizeError(ctx, c.Writer, parRequest, err)
		return
	}

	response, err := a.provider.NewPushedAuthorizeResponse(ctx, parRequest, new(models.Session))
	if err != nil {
		a.provider.WritePushedAuthorizeError(ctx, c.Writer, parRequest, err)
		return
	}

	a.provider.WritePushedAuthorizeResponse(ctx, c.Writer, parRequest, response)
}

func (a *Auth) RevokeHandler(c *gin.Context) {
	ctx := c.Request.Context()

	err := a.provider.NewRevocationRequest(ctx, c.Request)
	if err != nil {
		a.provider.WriteRevocationResponse(ctx, c.Writer, err)
		return
	}

	a.provider.WriteRevocationResponse(ctx, c.Writer, nil)
}

func (a *Auth) IntrospectionHandler(c *gin.Context) {
	ctx := c.Request.Context()

	ir, err := a.provider.NewIntrospectionRequest(ctx, c.Request, new(models.Session))
	if err != nil {
		a.provider.WriteIntrospectionError(ctx, c.Writer, err)
		return
	}

	a.provider.WriteIntrospectionResponse(ctx, c.Writer, ir)
}
