package routes

import (
	"net/http"

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
	r.GET("/stats", handlers.StatsPage)
	r.GET("/admin", handlers.AdminDashboardPage)
	r.GET("/admin/dashboard", handlers.AdminDashboardPage)
	r.GET("/admin/users", handlers.AdminUsersPage)
	r.GET("/admin/logs", handlers.AdminLogsPage)

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

		// Exercises (solo lectura para todos los usuarios autenticados)
		apiPrivate.GET("/exercises", exerciseHandler.GetExercise)
		apiPrivate.GET("/exercises/:id", exerciseHandler.GetExercise)

		// Exercises admin operations (solo administradores pueden crear/editar/eliminar)
		exercisesAdmin := apiPrivate.Group("/exercises")
		exercisesAdmin.Use(middleware.RequireRole("admin"))
		{
			exercisesAdmin.POST("", exerciseHandler.CreateExercise)
			exercisesAdmin.PUT("/:id", exerciseHandler.UpdateExercise)
			exercisesAdmin.DELETE("/:id", exerciseHandler.DeleteExercise)
		}

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

		// Admin API endpoints (require admin role)
		// El middleware RequireRole("admin") protege estas rutas
		adminAPI := apiPrivate.Group("/admin")
		adminAPI.Use(middleware.RequireRole("admin"))
		{
			// Obtener usuario por ID (usa el handler existente)
			adminAPI.GET("/users/:id", userHandler.GetUserByID)

			// Placeholders para endpoints que requieren handlers adicionales
			adminAPI.GET("/users", func(c *gin.Context) {
				c.JSON(http.StatusNotImplemented, gin.H{"error": "Endpoint not implemented. Use /api/users/:id"})
			})
			adminAPI.PUT("/users/:id/role", func(c *gin.Context) {
				c.JSON(http.StatusNotImplemented, gin.H{"error": "Handler not implemented yet"})
			})
			adminAPI.GET("/logs", func(c *gin.Context) {
				c.JSON(http.StatusNotImplemented, gin.H{"error": "Handler not implemented yet", "logs": []interface{}{}})
			})
			adminAPI.DELETE("/logs", func(c *gin.Context) {
				c.JSON(http.StatusNotImplemented, gin.H{"error": "Handler not implemented yet"})
			})
		}
	}
}
