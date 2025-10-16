package handlers

import (
	"net/http"

	"backend/dto"
	"backend/services"

	"github.com/gin-gonic/gin"
)

type RoutineHandler struct {
	service services.RoutineServiceInterface
}

func NewRoutineHandler(service services.RoutineServiceInterface) *RoutineHandler {
	return &RoutineHandler{service: service}
}

func (h *RoutineHandler) GetRoutine(c *gin.Context) {
	if id := c.Param("id"); id != "" {
		routine, err := h.service.GetRoutineByID(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Routine not found"})
			return
		}
		c.JSON(http.StatusOK, routine)
		return
	}

	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}
	name := c.Query("name")
	routines, err := h.service.GetRoutines(userID.(string), name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch routines"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"routines": routines})
}

func (h *RoutineHandler) CreateRoutine(c *gin.Context) {
	var req dto.RoutineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}
	req.UserID = userID.(string)

	routine, err := h.service.CreateRoutine(userID.(string), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, routine)
}

func (h *RoutineHandler) UpdateRoutine(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing routine ID"})
		return
	}

	var req dto.RoutineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	updated, err := h.service.UpdateRoutine(userID.(string), id, req)
	if err != nil {
		if err.Error() == "no autorizado: no es el owner de la rutina" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updated)
}

func (h *RoutineHandler) DeleteRoutine(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing routine ID"})
		return
	}

	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	err := h.service.DeleteRoutine(userID.(string), id)
	if err != nil {
		if err.Error() == "no autorizado: no es el owner de la rutina" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Routine deleted successfully"})
}
