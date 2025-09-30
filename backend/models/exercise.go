package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Exercise struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID      string             `bson:"user_id" json:"user_id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
	Category    string             `bson:"category" json:"category"`
	MuscleGroup string             `bson:"muscle_group" json:"muscle_group"`
	Difficulty  string             `bson:"difficulty" json:"difficulty"`
	MediaURL    string             `bson:"media_url,omitempty" json:"media_url,omitempty"`
	Steps       []string           `bson:"steps,omitempty" json:"steps,omitempty"`
	CreatedBy   primitive.ObjectID `bson:"created_by" json:"created_by"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}
