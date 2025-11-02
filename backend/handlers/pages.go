package handlers

import "github.com/gin-gonic/gin"

func LoginPage(c *gin.Context) {
	c.HTML(200, "layout", gin.H{
		"TemplateName": "login",
		"Title":        "Iniciar Sesión",
	})
}

func RegisterPage(c *gin.Context) {
	c.HTML(200, "layout", gin.H{
		"TemplateName": "register",
		"Title":        "Registrarse",
	})
}

func IndexPage(c *gin.Context) {
	c.HTML(200, "layout", gin.H{
		"TemplateName": "index",
		"Title":        "Inicio",
	})
}

func ExercisesPage(c *gin.Context) {
	c.HTML(200, "layout", gin.H{
		"TemplateName": "exercises",
		"Title":        "Ejercicios",
	})
}

func ExerciseCreatePage(c *gin.Context) {
	c.HTML(200, "layout", gin.H{
		"TemplateName": "exercise-create",
		"Title":        "Crear Ejercicio",
	})
}

func RoutinesPage(c *gin.Context) {
	c.HTML(200, "layout", gin.H{
		"TemplateName": "routines",
		"Title":        "Rutinas",
	})
}

func RoutineCreatePage(c *gin.Context) {
	c.HTML(200, "layout", gin.H{
		"TemplateName": "routine-create",
		"Title":        "Crear Rutina",
	})
}

func WorkoutsPage(c *gin.Context) {
	c.HTML(200, "layout", gin.H{
		"TemplateName": "workouts",
		"Title":        "Entrenamientos",
	})
}

func WorkoutCreatePage(c *gin.Context) {
	c.HTML(200, "layout", gin.H{
		"TemplateName": "workout-create",
		"Title":        "Registrar Entrenamiento",
	})
}

func ExerciseDetailPage(c *gin.Context) {
	c.HTML(200, "layout", gin.H{
		"TemplateName": "exercise-detail",
		"Title":        "Detalle del Ejercicio",
	})
}

func RoutineDetailPage(c *gin.Context) {
	c.HTML(200, "layout", gin.H{
		"TemplateName": "routine-detail",
		"Title":        "Detalle de la Rutina",
	})
}

func WorkoutDetailPage(c *gin.Context) {
	c.HTML(200, "layout", gin.H{
		"TemplateName": "workout-detail",
		"Title":        "Detalle del Entrenamiento",
	})
}

func ProfilePage(c *gin.Context) {
	c.HTML(200, "layout", gin.H{
		"TemplateName": "profile",
		"Title":        "Mi Perfil",
	})
}
