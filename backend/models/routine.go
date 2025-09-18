package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoutineEntry struct {
	ExerciseID primitive.ObjectID `bson:"exercise_id" json:"exercise_id"`
	Order      int                `bson:"order" json:"order"` // para ordenar ejercicios en la rutina
	Sets       int                `bson:"sets" json:"sets"`
	Reps       int                `bson:"reps" json:"reps"`
	Weight     float64            `bson:"weight,omitempty" json:"weight,omitempty"`
}

type Routine struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	OwnerID     primitive.ObjectID `bson:"owner_id" json:"owner_id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
	Entries     []RoutineEntry     `bson:"entries" json:"entries"`
	IsPublic    bool               `bson:"is_public" json:"is_public"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}
