package handlers

import (
	"backend/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type LogHandler struct {
	service services.LogServiceInterface
}

func NewLogHandler(service services.LogServiceInterface) *LogHandler {
	return &LogHandler{
		service: service,
	}
}

func (h *LogHandler) GetLogs(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "100")
	skipStr := c.DefaultQuery("skip", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 100
	}

	skip, err := strconv.Atoi(skipStr)
	if err != nil {
		skip = 0
	}

	logs, err := h.service.GetLogs(limit, skip)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener logs", "logs": []interface{}{}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"logs": logs})
}

func (h *LogHandler) DeleteLogs(c *gin.Context) {
	err := h.service.DeleteAllLogs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar logs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logs eliminados correctamente"})
}

