package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"backend/auth"
	"backend/dto"
	"backend/services"
)

type UserHandler struct {
	service services.UserServiceInterface
}

func NewUserHandler(service services.UserServiceInterface) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (handler *UserHandler) Register(c *gin.Context) {
	var request dto.RegisterRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := handler.service.Register(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (handler *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := handler.service.Login(req)
	if err != nil || user.ID == primitive.NilObjectID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "credenciales inválidas"})
		return
	}

	if !auth.CheckPasswordHash(req.Password, user.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "credenciales inválidas"})
		return
	}

	access, refresh, expiresIn, err := auth.GenerateToken(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudieron generar tokens"})
		return
	}

	c.JSON(http.StatusOK, dto.AuthResponse{
		AccessToken:  access,
		RefreshToken: refresh,
		ExpiresIn:    expiresIn,
		User:         user,
	})
}
