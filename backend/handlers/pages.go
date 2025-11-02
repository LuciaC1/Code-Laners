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
func ExerciseEditPage(c *gin.Context) {
	c.HTML(200, "layout", gin.H{
		"TemplateName": "exercise-edit",
		"Title":        "Editar Ejercicio",
	})
}
func RoutineDetailPage(c *gin.Context) {
	c.HTML(200, "layout", gin.H{
		"TemplateName": "routine-detail",
		"Title":        "Detalle de la Rutina",
	})
}
func RoutineEditPage(c *gin.Context) {
	c.HTML(200, "layout", gin.H{
		"TemplateName": "routine-edit",
		"Title":        "Editar Rutina",
	})
}
func WorkoutDetailPage(c *gin.Context) {
	c.HTML(200, "layout", gin.H{
		"TemplateName": "workout-detail",
		"Title":        "Detalle del Entrenamiento",
	})
}
func WorkoutEditPage(c *gin.Context) {
	c.HTML(200, "layout", gin.H{
		"TemplateName": "workout-edit",
		"Title":        "Editar Entrenamiento",
	})
}
func ProfilePage(c *gin.Context) {
	c.HTML(200, "layout", gin.H{
		"TemplateName": "profile",
		"Title":        "Mi Perfil",
	})
}
func StatsPage(c *gin.Context) {
	c.HTML(200, "layout", gin.H{
		"TemplateName": "stats",
		"Title":        "Mis Estadísticas",
	})
}
func AdminDashboardPage(c *gin.Context) {
	c.HTML(200, "layout", gin.H{
		"TemplateName": "admin-dashboard",
		"Title":        "Panel de Administración",
	})
}
func AdminUsersPage(c *gin.Context) {
	c.HTML(200, "layout", gin.H{
		"TemplateName": "admin-users",
		"Title":        "Gestión de Usuarios",
	})
}
func AdminLogsPage(c *gin.Context) {
	c.HTML(200, "layout", gin.H{
		"TemplateName": "admin-logs",
		"Title":        "Logs del Sistema",
	})
}
