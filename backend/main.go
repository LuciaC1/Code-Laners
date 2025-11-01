package main

import (
	"log"
	"net/http"

	"backend/database"
	"backend/handlers"
	"backend/repositories"
	"backend/routes"
	"backend/services"

	"github.com/gin-gonic/gin"
)

func main() {
	db := database.NewMongoDB()
	if err := db.Connect(); err != nil {
		log.Fatal("Error conectando a MongoDB:", err)
	}
	defer db.Disconnect()

	exRepo := repositories.NewExerciseRepository(db)
	rtRepo := repositories.NewRoutineRepository(db)
	wkRepo := repositories.NewWorkoutRepository(db)

	exSvc := services.NewExerciseService(exRepo)
	rtSvc := services.NewRoutineService(rtRepo, exRepo)
	wkSvc := services.NewWorkoutService(wkRepo)

	// Handlers (según carpeta backend/handlers)
	exHandler := handlers.NewExerciseHandler(exSvc)
	rtHandler := handlers.NewRoutineHandler(rtSvc)
	wkHandler := handlers.NewWorkoutHandler(wkSvc)

	r := gin.Default()

	// Templates y static (rutas relativas a la raíz del repo)
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")

	// Ruta raíz que renderiza index.html (ajusta nombre si usas layout diferente)
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	// Registrar rutas con handlers instanciados (usar nil si no hay userHandler)
	routes.SetupRoutes(r, nil, exHandler, rtHandler, wkHandler)

	log.Println("Servidor iniciado en http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Error iniciando servidor:", err)
	}
}
