package main

import (
	"backend/auth"
	"backend/database"
	"backend/handlers"
	"backend/models"
	"backend/repositories"
	"backend/routes"
	"backend/services"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func initializeDefaultAdmin(db database.DB) {
	userRepo := repositories.NewUserRepository(db)
	
	_, err := userRepo.GetUserByEmail("cuentaadmin@admin.com")
	
	if err == mongo.ErrNoDocuments {
		hash, err := auth.HashPassword("12345678")
		if err != nil {
			log.Printf("Error hasheando contraseña del admin: %v", err)
			return
		}

		now := time.Now()
		adminUser := models.User{
			ID:           primitive.NewObjectID(),
			Name:         "Admin",
			Email:        "cuentaadmin@admin.com",
			PasswordHash: string(hash),
			Role:         models.RoleAdmin,
			DateOfBirth:  time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC), 
			CreatedAt:    now,
			UpdatedAt:    now,
		}
		
		_, err = userRepo.CreateUser(adminUser)
		if err != nil {
			log.Printf("Error creando usuario admin por defecto: %v", err)
			return
		}
		
		log.Println("Usuario admin por defecto creado exitosamente")
	} else if err != nil {
		log.Printf("Error verificando usuario admin: %v", err)
	} else {
		log.Println("Usuario admin por defecto ya existe")
	}
}

func main() {
	db := database.NewMongoDB()
	if err := db.Connect(); err != nil {
		log.Fatal("Error conectando a MongoDB:", err)
	}
	defer db.Disconnect()
	
	initializeDefaultAdmin(db)
	
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	r.Static("/static", "./static")

	statsRepo := repositories.NewStatsRepository(db)
	statsService := services.NewStatsService(statsRepo, repositories.NewRoutineRepository(db), repositories.NewUserRepository(db), repositories.NewExerciseRepository(db))
	statsHandler := handlers.NewStatsHandler(statsService)

	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	
	logRepo := repositories.NewLogRepository(db)
	logService := services.NewLogService(logRepo)
	logHandler := handlers.NewLogHandler(logService)

	userHandler := handlers.NewUserHandler(userService, logService)
	exerciseHandler := handlers.NewExerciseHandler(
		services.NewExerciseService(repositories.NewExerciseRepository(db)),
		logService,
		userService,
	)
	routineHandler := handlers.NewRoutineHandler(
		services.NewRoutineService(repositories.NewRoutineRepository(db), repositories.NewExerciseRepository(db), userRepo),
		logService,
		userService,
	)
	workoutHandler := handlers.NewWorkoutHandler(
		services.NewWorkoutService(repositories.NewWorkoutRepository(db), repositories.NewRoutineRepository(db)),
		logService,
		userService,
	)

	routes.SetupRoutes(
		r,
		userHandler,
		exerciseHandler,
		routineHandler,
		workoutHandler,
		statsHandler,
		logHandler,
	)
	log.Println("Servidor iniciado en http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Error iniciando servidor:", err)
	}
}
