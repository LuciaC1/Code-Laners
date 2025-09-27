package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ExercisePerformance struct {
	ExerciseID primitive.ObjectID `bson:"exercise_id" json:"exercise_id"`
	Sets       int                `bson:"sets" json:"sets"`
	Reps       int                `bson:"reps" json:"reps"`
	Weight     float64            `bson:"weight,omitempty" json:"weight,omitempty"`
}

type Workout struct {
	ID                 primitive.ObjectID    `bson:"_id,omitempty" json:"id"`
	UserID             primitive.ObjectID    `bson:"user_id" json:"user_id"`
	RoutineID          *primitive.ObjectID   `bson:"routine_id,omitempty" json:"routine_id,omitempty"`
	PerformedExercises []ExercisePerformance `bson:"performed_exercises" json:"performed_exercises"`
	CompletedAt        time.Time             `bson:"completed_at" json:"completed_at"`
	UpdatedAt          time.Time             `bson:"updated_at" json:"updated_at"`
	DurationMinutes    int                   `bson:"duration_minutes,omitempty" json:"duration_minutes,omitempty"`
	Notes              string                `bson:"notes,omitempty" json:"notes,omitempty"`
}
