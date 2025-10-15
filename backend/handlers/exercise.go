package handlers

import (
	"backend/dto"
	"backend/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ExerciseHandler struct {
	service services.ExerciseInterface
}

func NewExerciseHandler(service services.ExerciseInterface) *ExerciseHandler {
	return &ExerciseHandler{service: service}
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
	var exerciseReq dto.ExerciseRequest
	if err := c.ShouldBindJSON(&exerciseReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}
	exerciseReq.UserID = userID.(string)

	exercise, err := h.service.CreateExercise(exerciseReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create exercise"})
		return
	}

	c.JSON(http.StatusCreated, exercise)
}

func (h *ExerciseHandler) UpdateExercise(c *gin.Context) {
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

	userID, _ := c.Get("user_id")
	exerciseReq.UserID = userID.(string)

	exercise, err := h.service.UpdateExercise(id, exerciseReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update exercise"})
		return
	}

	c.JSON(http.StatusOK, exercise)
}

func (h *ExerciseHandler) DeleteExercise(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing exercise ID"})
		return
	}

	userRole, _ := c.Get("user_role")
	role := userRole.(string)

	err := h.service.DeleteExercise(id, role)
	if err != nil {
		if err.Error() == "forbidden: only admins can delete exercises" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete exercise"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Exercise deleted successfully"})
}
