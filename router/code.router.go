package router

func (r *router) CodeRouter() {
	codeCroup := r.Server.Group("/code")

	{
		codeCroup.POST("/verif/:id", r.StoreRequest.VerifCode)
		codeCroup.POST("/restart/:name/:table", r.StoreRequest.RestartCode)
	}
}
