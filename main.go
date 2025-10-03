package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dylEasydev/go-oauth2-easyclass/controller"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv"
)

func main() {

	server := gin.Default()

	port := os.Getenv("PORT")

	if err := server.Run(":" + port); err != nil {
		log.Fatal("Erreur du démarrage du serveur", err)
	}

	server.GET("/", controller.IndexHanler)

	fmt.Printf("Serveur démarre à l'adresse http://localhost:%s", port)
}
