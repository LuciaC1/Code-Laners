package dto

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct { //Revisar este struct, esta copiado de models.User
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name         string             `bson:"name" json:"name"`
	Email        string             `bson:"email" json:"email"`
	PasswordHash string             `bson:"password_hash" json:"-"`
	Role         string             `bson:"role" json:"role"`
	DateOfBirth  time.Time          `bson:"date_of_birth" json:"date_of_birth"`
	Weight       float64            `bson:"weight,omitempty" json:"weight,omitempty"`
	Height       float64            `bson:"height,omitempty" json:"height,omitempty"`
	Level        string             `bson:"level,omitempty" json:"level,omitempty"`
	Goals        []string           `bson:"goals,omitempty" json:"goals,omitempty"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
}

type RegisterRequest struct {
	Name        string   `json:"name" binding:"required,min=2,max=100"`
	Email       string   `json:"email" binding:"required,email"`
	Password    string   `json:"password" binding:"required,min=6"`
	DateOfBirth string   `json:"date_of_birth" binding:"required"` // ISO date string - parsear en servicio
	Weight      *float64 `json:"weight,omitempty"`
	Height      *float64 `json:"height,omitempty"`
	Level       string   `json:"level,omitempty"`
	Goals       []string `json:"goals,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
type UpdateUserRequest struct {
	Name   string   `json:"name,omitempty"`
	Email  string   `json:"email,omitempty"` // validar en servicio
	Weight float64  `json:"weight,omitempty"`
	Height float64  `json:"height,omitempty"`
	Level  string   `json:"level,omitempty"`
	Goals  []string `json:"goals,omitempty"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}
