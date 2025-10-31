package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"backend/auth"
	"backend/database"
	"backend/dto"
	"backend/models"
	"backend/repositories"
	"backend/services"
	"time"
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

	access, refresh, expiresIn, err := auth.GenerateToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudieron generar tokens"})
		return
	}

	db := database.NewMongoDB()
	refreshRepo := repositories.NewRefreshTokenRepository(db)
	rt := models.RefreshToken{
		UserID:    user.ID,
		Token:     refresh,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		Revoked:   false,
	}
	_, _ = refreshRepo.Save(rt)

	c.JSON(http.StatusOK, dto.AuthResponse{
		AccessToken:  access,
		RefreshToken: refresh,
		ExpiresIn:    expiresIn,
		User:         user,
	})
}

func (handler *UserHandler) Refresh(c *gin.Context) {
	var req dto.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.RefreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "refreshToken es requerido"})
		return
	}

	claims, err := auth.ValidateToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token inválido o expirado"})
		return
	}

	db := database.NewMongoDB()
	refreshRepo := repositories.NewRefreshTokenRepository(db)
	saved, err := refreshRepo.GetByToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token inválido o revocado"})
		return
	}
	if saved.Revoked || saved.ExpiresAt.Before(time.Now()) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token inválido o revocado"})
		return
	}

	access, expiresIn, err := auth.GenerateAccessTokenFromStrings(claims.UserID, claims.Email, claims.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo generar access token"})
		return
	}

	c.JSON(http.StatusOK, dto.RefreshResponse{
		AccessToken: access,
		ExpiresIn:   expiresIn,
	})
}

func (handler *UserHandler) Logout(c *gin.Context) {
	var req dto.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.RefreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "refreshToken es requerido"})
		return
	}

	_, err := auth.ValidateToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "refresh token inválido"})
		return
	}

	db := database.NewMongoDB()
	refreshRepo := repositories.NewRefreshTokenRepository(db)
	_, err = refreshRepo.Revoke(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo revocar el token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logout exitoso"})
}
