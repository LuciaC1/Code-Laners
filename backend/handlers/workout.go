package handlers
import (
	"net/http"
	"time"
	"backend/dto"
	"backend/services"
	"github.com/gin-gonic/gin"
)
type WorkoutHandler struct {
	service services.WorkoutServiceInterface
}
func NewWorkoutHandler(s services.WorkoutServiceInterface) *WorkoutHandler {
	return &WorkoutHandler{service: s}
}
func (h *WorkoutHandler) GetWorkouts(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}
	workouts, err := h.service.GetWorkouts(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"workouts": workouts})
}
func (h *WorkoutHandler) CreateWorkout(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}
	var workout dto.WorkoutDTO
	if err := c.ShouldBindJSON(&workout); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	workout.UserID = userID.(string)
	workout.CompletedAt = time.Now()
	workout.UpdatedAt = time.Now()
	id, err := h.service.CreateWorkout(workout)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}
func (h *WorkoutHandler) GetWorkoutByID(c *gin.Context) {
	workoutID := c.Param("id")
	if workoutID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de workout requerido"})
		return
	}
	workout, err := h.service.GetWorkoutByID(workoutID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workout no encontrado"})
		return
	}
	c.JSON(http.StatusOK, workout)
}
func (h *WorkoutHandler) UpdateWorkout(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}
	workoutID := c.Param("id")
	if workoutID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de workout requerido"})
		return
	}
	existingWorkout, err := h.service.GetWorkoutByID(workoutID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workout no encontrado"})
		return
	}
	if existingWorkout.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "No tienes permiso para modificar este workout"})
		return
	}
	var workout dto.WorkoutDTO
	if err := c.ShouldBindJSON(&workout); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	workout.ID = existingWorkout.ID
	workout.UserID = userID.(string)
	workout.UpdatedAt = time.Now()
	if err := h.service.UpdateWorkout(workout); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Workout actualizado exitosamente"})
}
func (h *WorkoutHandler) DeleteWorkout(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}
	workoutID := c.Param("id")
	if workoutID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de workout requerido"})
		return
	}
	existingWorkout, err := h.service.GetWorkoutByID(workoutID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workout no encontrado"})
		return
	}
	if existingWorkout.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "No tienes permiso para eliminar este workout"})
		return
	}
	if err := h.service.DeleteWorkout(workoutID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Workout eliminado exitosamente"})
}
