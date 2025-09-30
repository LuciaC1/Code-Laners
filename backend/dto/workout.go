package dto

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Workout struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID            string             `bson:"user_id" json:"user_id"`
	RoutineID         string             `bson:"routine_id,omitempty" json:"routine_id,omitempty"`
	CompletedAt       time.Time          `bson:"completed_at" json:"completed_at"`
	UpdatedAt         time.Time          `bson:"updated_at" json:"updated_at"`
	DurationMinutes   int                `bson:"duration_minutes,omitempty" json:"duration_minutes,omitempty"`
	Notes             string             `bson:"notes,omitempty" json:"notes,omitempty"`
	EstimatedCalories int                `bson:"estimated_calories,omitempty" json:"estimated_calories,omitempty"`
}
