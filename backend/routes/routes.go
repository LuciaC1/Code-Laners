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

	// Service HTML pages (may require authentication)
	r.GET("/exercises", handlers.ExercisesPage)
	r.GET("/exercises/create", handlers.ExerciseCreatePage)
	r.GET("/exercises/:id", handlers.ExerciseDetailPage)
	r.GET("/routines", handlers.RoutinesPage)
	r.GET("/routines/create", handlers.RoutineCreatePage)
	r.GET("/routines/:id", handlers.RoutineDetailPage)
	r.GET("/workouts", handlers.WorkoutsPage)
	r.GET("/workouts/create", handlers.WorkoutCreatePage)
	r.GET("/workouts/:id", handlers.WorkoutDetailPage)
	r.GET("/profile", handlers.ProfilePage)

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
		// User Profile
		apiPrivate.GET("/me", userHandler.GetMe)
		apiPrivate.PUT("/me", userHandler.UpdateMe)
		apiPrivate.PUT("/me/password", userHandler.ChangePassword)

		// Exercises
		apiPrivate.GET("/exercises", exerciseHandler.GetExercise)
		apiPrivate.GET("/exercises/:id", exerciseHandler.GetExercise)
		apiPrivate.POST("/exercises", exerciseHandler.CreateExercise)
		apiPrivate.PUT("/exercises/:id", exerciseHandler.UpdateExercise)
		apiPrivate.DELETE("/exercises/:id", exerciseHandler.DeleteExercise)

		// Routines
		apiPrivate.GET("/routines", routineHandler.GetRoutines)
		apiPrivate.GET("/routines/:id", routineHandler.GetRoutineByID)
		apiPrivate.POST("/routines", routineHandler.CreateRoutine)
		apiPrivate.PUT("/routines/:id", routineHandler.UpdateRoutine)
		apiPrivate.DELETE("/routines/:id", routineHandler.DeleteRoutine)

		// Workouts
		apiPrivate.GET("/workouts", workoutHandler.GetWorkouts)
		apiPrivate.GET("/workouts/:id", workoutHandler.GetWorkoutByID)
		apiPrivate.POST("/workouts", workoutHandler.CreateWorkout)
		apiPrivate.PUT("/workouts/:id", workoutHandler.UpdateWorkout)
		apiPrivate.DELETE("/workouts/:id", workoutHandler.DeleteWorkout)
	}
}
