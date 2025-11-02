package main

import (
	"backend/database"
	"backend/handlers"
	"backend/repositories"
	"backend/routes"
	"backend/services"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	db := database.NewMongoDB()
	if err := db.Connect(); err != nil {
		log.Fatal("Error conectando a MongoDB:", err)
	}
	defer db.Disconnect()
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	// Archivos estáticos
	r.Static("/static", "./static")
	// Registrar rutas (delegado a routes.SetupRoutes)
	routes.SetupRoutes(
		r,
		handlers.NewUserHandler(services.NewUserService(repositories.NewUserRepository(db))),
		handlers.NewExerciseHandler(services.NewExerciseService(repositories.NewExerciseRepository(db))),
		handlers.NewRoutineHandler(services.NewRoutineService(repositories.NewRoutineRepository(db), repositories.NewExerciseRepository(db), repositories.NewUserRepository(db))),
		handlers.NewWorkoutHandler(services.NewWorkoutService(repositories.NewWorkoutRepository(db), repositories.NewRoutineRepository(db))),
	)
	log.Println("Servidor iniciado en http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Error iniciando servidor:", err)
	}
}
