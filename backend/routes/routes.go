package routes

import (
	"backend/handlers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, userHandler *handlers.UserHandler, exerciseHandler *handlers.ExerciseHandler, routineHandler *handlers.RoutineHandler, workoutHandler *handlers.WorkoutHandler) {
	api := r.Group("/api")

	// Auth routes
	auth := api.Group("/auth")
	{
		auth.POST("/register", userHandler.Register)
		auth.POST("/login", userHandler.Login)
		auth.POST("/refresh", userHandler.Refresh)
		auth.POST("/logout", userHandler.Logout)
	}

	me := api.Group("/me")
	me.Use(middleware.AuthMiddleware())
	{
		me.GET("/", userHandler.GetMe)
	}

	exercises := api.Group("/exercises")
	{
		exercises.GET("", exerciseHandler.GetExercise)
		exercises.POST("", middleware.AuthMiddleware(), exerciseHandler.CreateExercise)
		exercises.PUT("/:id", middleware.AuthMiddleware(), exerciseHandler.UpdateExercise)
		exercises.DELETE("/:id", middleware.AuthMiddleware(), exerciseHandler.DeleteExercise)
	}

	routines := api.Group("/routines")
	{
		routines.GET("", routineHandler.GetRoutines)
		routines.GET("/:id", routineHandler.GetRoutineByID)

		routines.POST("", middleware.AuthMiddleware(), routineHandler.CreateRoutine)
		routines.PUT("/:id", middleware.AuthMiddleware(), routineHandler.UpdateRoutine)
		routines.DELETE("/:id", middleware.AuthMiddleware(), routineHandler.DeleteRoutine)
	}

	workouts := api.Group("/workouts")
	{
		workouts.GET("", workoutHandler.GetWorkouts)
		workouts.GET("/:id", workoutHandler.GetWorkoutByID)

		workouts.POST("", middleware.AuthMiddleware(), workoutHandler.CreateWorkout)
		workouts.PUT("/:id", middleware.AuthMiddleware(), workoutHandler.UpdateWorkout)
		workouts.DELETE("/:id", middleware.AuthMiddleware(), workoutHandler.DeleteWorkout)
	}
}
