package routes

import (
	"backend/handlers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, userHandler *handlers.UserHandler, exerciseHandler *handlers.ExerciseHandler, routineHandler *handlers.RoutineHandler, workoutHandler *handlers.WorkoutHandler) {
	// Public HTML pages
	r.GET("/", handlers.IndexPage)
	r.GET("/login", handlers.LoginPage)
	r.GET("/register", handlers.RegisterPage)
	r.GET("/products", handlers.ProductsPage)

	// API pública
	api := r.Group("/api")
	{
		api.POST("/register", userHandler.Register)
		api.POST("/login", userHandler.Login)
	}

	// API privada (requiere autenticación)
	apiPrivate := r.Group("/api")
	apiPrivate.Use(middleware.AuthMiddleware())
	{
		apiPrivate.GET("/exercises", exerciseHandler.GetExercise)
		apiPrivate.POST("/exercises", exerciseHandler.CreateExercise)
		apiPrivate.GET("/routines", routineHandler.GetRoutines)
		apiPrivate.POST("/routines", routineHandler.CreateRoutine)
		apiPrivate.GET("/workouts", workoutHandler.GetWorkouts)
		apiPrivate.POST("/workouts", workoutHandler.CreateWorkout)
	}
}
