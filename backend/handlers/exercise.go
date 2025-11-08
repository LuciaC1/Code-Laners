package handlers

import (
	"backend/dto"
	"backend/models"
	"backend/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ExerciseHandler struct {
	service    services.ExerciseInterface
	logService services.LogServiceInterface
	userService services.UserServiceInterface
}

func NewExerciseHandler(service services.ExerciseInterface, logService services.LogServiceInterface, userService services.UserServiceInterface) *ExerciseHandler {
	return &ExerciseHandler{
		service:     service,
		logService:  logService,
		userService: userService,
	}
}
func (h *ExerciseHandler) GetExercise(c *gin.Context) {
	if id := c.Param("id"); id != "" {
		exercise, err := h.service.GetExerciseByID(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Exercise not found"})
			return
		}
		c.JSON(http.StatusOK, exercise)
		return
	}
	var search dto.ExerciseSearch
	if err := c.ShouldBindQuery(&search); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	exercises, err := h.service.GetExercises(search.Name, search.Category, search.MuscleGroup)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch exercises"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"exercises": exercises})
}
func (h *ExerciseHandler) CreateExercise(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}
	var exerciseReq dto.ExerciseRequest
	if err := c.ShouldBindJSON(&exerciseReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	exerciseReq.UserID = userID.(string)
	exercise, err := h.service.CreateExercise(exerciseReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create exercise"})
		return
	}

	if h.logService != nil {
		userIDObj, _ := primitive.ObjectIDFromHex(userID.(string))
		user, _ := h.userService.GetUserByID(userID.(string))
		_ = h.logService.CreateLog(
			models.LogLevelSuccess,
			models.LogTypeExercise,
			"Ejercicio creado",
			&userIDObj,
			user.Name,
			"Ejercicio: "+exercise.Name,
		)
	}
	c.JSON(http.StatusCreated, exercise)
}
func (h *ExerciseHandler) UpdateExercise(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing exercise ID"})
		return
	}
	var exerciseReq dto.ExerciseRequest
	if err := c.ShouldBindJSON(&exerciseReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	exerciseReq.UserID = userID.(string)
	exercise, err := h.service.UpdateExercise(id, exerciseReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update exercise"})
		return
	}

	if h.logService != nil {
		userIDObj, _ := primitive.ObjectIDFromHex(userID.(string))
		user, _ := h.userService.GetUserByID(userID.(string))
		_ = h.logService.CreateLog(
			models.LogLevelInfo,
			models.LogTypeExercise,
			"Ejercicio actualizado",
			&userIDObj,
			user.Name,
			"Ejercicio: "+exercise.Name,
		)
	}
	c.JSON(http.StatusOK, exercise)
}
func (h *ExerciseHandler) DeleteExercise(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing exercise ID"})
		return
	}

	exercise, err := h.service.GetExerciseByID(id)
	if err == nil && h.logService != nil {
		userIDObj, _ := primitive.ObjectIDFromHex(userID.(string))
		user, _ := h.userService.GetUserByID(userID.(string))
		_ = h.logService.CreateLog(
			models.LogLevelWarning,
			models.LogTypeExercise,
			"Ejercicio eliminado",
			&userIDObj,
			user.Name,
			"Ejercicio: "+exercise.Name,
		)
	}
	err = h.service.DeleteExercise(userID.(string), id)
	if err != nil {
		if err.Error() == "unauthorized: cannot delete exercise you do not own" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Exercise deleted successfully"})
}
