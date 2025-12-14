package main

import (
	"log"
	"os"

	"github.com/dylEasydev/go-oauth2-easyclass/db"
	"github.com/dylEasydev/go-oauth2-easyclass/middleware"
	"github.com/dylEasydev/go-oauth2-easyclass/router"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	//chargement du fichier env
	_ = godotenv.Load()

	//initailisation du serveur
	server := gin.Default()

	//intialisation de la BD
	store := db.New()

	port := os.Getenv("PORT")

	log.Printf("Serveur démarre à l'adresse https://localhost:%s", port)

	server.Use(middleware.ErrorHandler())
	router := router.NewRouter(server, store)
	router.IndexRouter()
	router.OIDCRouter()
	router.JWKRouter()
	router.SignRouter()
	router.CodeRouter()

	//démarrage du serveur https
	if err := server.RunTLS(":"+port, "./key/server.pem", "./key/server.key"); err != nil {
		log.Fatal("Erreur du démarrage du serveur", err)
	}

}
