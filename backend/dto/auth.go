package dto

import "time"

type RegisterDTO struct {
	Name        string   `json:"name" binding:"required,min=2,max=100"`
	Email       string   `json:"email" binding:"required,email"`
	Password    string   `json:"password" binding:"required,min=6"`
	DateOfBirth string   `json:"date_of_birth" binding:"required"` // ISO date string - parsear en servicio
	Weight      *float64 `json:"weight,omitempty"`
	Height      *float64 `json:"height,omitempty"`
	Level       string   `json:"level,omitempty"`
	Goals       []string `json:"goals,omitempty"`
}

type LoginDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type TokenResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}
