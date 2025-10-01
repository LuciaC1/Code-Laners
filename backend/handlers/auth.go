package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"backend/auth"
	"backend/dto"
	"backend/services"
)

func Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if u, err := services.(req.Email); err == nil && u.ID != primitive.NilObjectID {
		c.JSON(http.StatusConflict, gin.H{"error": "usuario ya existe"})
		return
	}

	// 2) preparar usuario (hash password, parse date, asignar role)
	pwHash, err := auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al hashear contraseña"})
		return
	}
	dob, err := time.Parse(time.RFC3339, req.DateOfBirth)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "date_of_birth debe ser ISO (RFC3339)"})
		return
	}

	newUser := dto.User{
		ID:           primitive.NewObjectID(),
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: pwHash,
		Role:         "user",
		DateOfBirth:  dob,
		Weight:       0,
		Height:       0,
		Level:        req.Level,
		Goals:        req.Goals,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	if req.Weight != nil {
		newUser.Weight = *req.Weight
	}
	if req.Height != nil {
		newUser.Height = *req.Height
	}

	// 3) pedir al service que cree el usuario (service hace la BD)
	created, err := services.CreateUser(newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo crear usuario"})
		return
	}

	// 4) generar tokens
	access, refresh, expiresIn, err := auth.GenerateTokens(created.ID.Hex(), created.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudieron generar tokens"})
		return
	}

	// 5) opcional: guardar refresh token (en service / BD) para revocación/rotación
	_ = service.SaveRefreshToken(created.ID, refresh)

	// 6) devolver AuthResponse (user DTO ya oculta password por tag `json:"-"`)
	c.JSON(http.StatusCreated, dto.AuthResponse{
		AccessToken:  access,
		RefreshToken: refresh,
		ExpiresIn:    expiresIn,
		User:         created,
	})
}

// Login handler: usar service.GetUser para obtener usuario y autenticar
func Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := service.GetUser(req.Email)
	if err != nil || user.ID == primitive.NilObjectID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "credenciales inválidas"})
		return
	}

	if !auth.CheckPasswordHash(req.Password, user.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "credenciales inválidas"})
		return
	}

	access, refresh, expiresIn, err := auth.GenerateTokens(user.ID.Hex(), user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudieron generar tokens"})
		return
	}

	_ = service.SaveRefreshToken(user.ID, refresh)

	c.JSON(http.StatusOK, dto.AuthResponse{
		AccessToken:  access,
		RefreshToken: refresh,
		ExpiresIn:    expiresIn,
		User:         user,
	})
}
