package main

import (
	"backend/database"
	"backend/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := database.NewMongoDB().Connect(); err != nil {
		log.Fatal("Error conectando a MongoDB:", err)
	}
	defer database.NewMongoDB().Disconnect()

	r := gin.Default()

	r.LoadHTMLGlob("templates/*")

	r.Static("/static", "./static")

	routes.SetupRoutes(r, nil, nil, nil, nil)

	log.Println("Servidor iniciado en http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Error iniciando servidor:", err)
	}
}
