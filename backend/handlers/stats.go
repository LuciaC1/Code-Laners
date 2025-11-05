package handlers

import (
	"backend/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type StatsHandler struct {
	service services.StatsServiceInterface
}

func NewStatsHandler(s services.StatsServiceInterface) *StatsHandler {
	return &StatsHandler{service: s}
}

func (h *StatsHandler) GetUserStats(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}
	stats, err := h.service.GetUserStats(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}

func (h *StatsHandler) GetAdminStats(c *gin.Context) {
	stats, err := h.service.GetAdminStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}
