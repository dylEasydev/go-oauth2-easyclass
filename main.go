package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dylEasydev/go-oauth2-easyclass/db"
	"github.com/dylEasydev/go-oauth2-easyclass/router"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv"
)

func main() {

	server := gin.Default()
	store := db.New()

	port := os.Getenv("PORT")

	if err := server.RunTLS(":"+port, "./key/server.pem", "./key/server.key"); err != nil {
		log.Fatal("Erreur du démarrage du serveur", err)
	}

	router := router.Router{Server: server, Store: store}
	router.IndexRouter()
	router.OIDCRouter()
	router.JWKRouter()

	fmt.Printf("Serveur démarre à l'adresse http://localhost:%s", port)
}
