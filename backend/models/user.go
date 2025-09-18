package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name         string             `bson:"name" json:"name"`
	Email        string             `bson:"email" json:"email"`
	PasswordHash string             `bson:"password_hash" json:"-"`
	Role         Role               `bson:"role" json:"role"`
	DateOfBirth  time.Time          `bson:"date_of_birth" json:"date_of_birth"`
	Weight       float64            `bson:"weight,omitempty" json:"weight,omitempty"`
	Height       float64            `bson:"height,omitempty" json:"height,omitempty"`
	Level        string             `bson:"level,omitempty" json:"level,omitempty"`
	Goals        []string           `bson:"goals,omitempty" json:"goals,omitempty"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
}
