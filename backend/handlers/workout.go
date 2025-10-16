package handlers

import (
	"net/http"
	"time"

	"backend/dto"
	"backend/services"

	"github.com/gin-gonic/gin"
)

func GetWorkout(c *gin.Context, workoutService services.WorkoutServiceInterface) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	workouts, err := workoutService.GetWorkouts(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, workouts)
}

func CreateWorkout(c *gin.Context, workoutService services.WorkoutServiceInterface) {
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

	id, err := workoutService.CreateWorkout(workout)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func UpdateWorkout(c *gin.Context, workoutService services.WorkoutServiceInterface) {
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

	// Verify the workout exists and belongs to the user
	existingWorkout, err := workoutService.GetWorkoutByID(workoutID)
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

	if err := workoutService.UpdateWorkout(workout); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Workout actualizado exitosamente"})
}

func DeleteWorkout(c *gin.Context, workoutService services.WorkoutServiceInterface) {
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

	// Verify the workout exists and belongs to the user
	existingWorkout, err := workoutService.GetWorkoutByID(workoutID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workout no encontrado"})
		return
	}

	if existingWorkout.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "No tienes permiso para eliminar este workout"})
		return
	}

	if err := workoutService.DeleteWorkout(workoutID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Workout eliminado exitosamente"})
}
