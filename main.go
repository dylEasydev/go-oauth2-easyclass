package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dylEasydev/go-oauth2-easyclass/db"
	"github.com/dylEasydev/go-oauth2-easyclass/middleware"
	"github.com/dylEasydev/go-oauth2-easyclass/router"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	_ = godotenv.Load()

	server := gin.Default()
	store := db.New()

	port := os.Getenv("PORT")

	if err := server.RunTLS(":"+port, "./key/server.pem", "./key/server.key"); err != nil {
		log.Fatal("Erreur du démarrage du serveur", err)
	}

	server.Use(middleware.ErrorHandler())
	router := router.NewRouter(server, store)
	router.IndexRouter()
	router.OIDCRouter()
	router.JWKRouter()
	router.SignRouter()
	router.CodeRouter()

	fmt.Printf("Serveur démarre à l'adresse https://localhost:%s", port)
}
