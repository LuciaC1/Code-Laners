package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LogLevel string

const (
	LogLevelInfo    LogLevel = "info"
	LogLevelWarning LogLevel = "warning"
	LogLevelError   LogLevel = "error"
	LogLevelSuccess LogLevel = "success"
)

type LogType string

const (
	LogTypeAuth     LogType = "auth"
	LogTypeUser     LogType = "user"
	LogTypeExercise LogType = "exercise"
	LogTypeRoutine  LogType = "routine"
	LogTypeWorkout  LogType = "workout"
	LogTypeGeneral  LogType = "general"
)

type Log struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Timestamp time.Time          `bson:"timestamp" json:"timestamp"`
	Level     LogLevel           `bson:"level" json:"level"`
	Type      LogType            `bson:"type" json:"type"`
	Action    string             `bson:"action" json:"action"`
	UserID    *primitive.ObjectID `bson:"user_id,omitempty" json:"user_id,omitempty"`
	UserName  string             `bson:"user_name,omitempty" json:"user_name,omitempty"`
	Details   string             `bson:"details,omitempty" json:"details,omitempty"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

