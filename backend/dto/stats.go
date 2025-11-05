package dto

import "time"

type WeeklyCount struct {
	Label string `json:"label"`
	Count int64  `json:"count"`
}

type RoutineCount struct {
	RoutineID   string `json:"routine_id,omitempty"`
	RoutineName string `json:"routine_name"`
	Count       int64  `json:"count"`
}

type ProgressPoint struct {
	Date     time.Time `json:"date"`
	Calories int       `json:"calories"`
	Duration int       `json:"duration"`
}

type RecentWorkout struct {
	ID                string    `json:"id"`
	CompletedAt       time.Time `json:"completed_at"`
	RoutineName       string    `json:"routine_name,omitempty"`
	DurationMinutes   int       `json:"duration_minutes,omitempty"`
	EstimatedCalories int       `json:"estimated_calories,omitempty"`
}

type UserStatsResponse struct {
	Summary struct {
		TotalWorkouts int64 `json:"totalWorkouts"`
		TotalRoutines int64 `json:"totalRoutines"`
		TotalMinutes  int64 `json:"totalMinutes"`
		TotalCalories int64 `json:"totalCalories"`
	} `json:"summary"`
	Frequency            []WeeklyCount   `json:"frequency"`
	RoutinesDistribution []RoutineCount  `json:"routinesDistribution"`
	Progress             []ProgressPoint `json:"progress"`
	RecentWorkouts       []RecentWorkout `json:"recentWorkouts"`
}

type AdminTotals struct {
	TotalUsers     int64 `json:"totalUsers"`
	TotalExercises int64 `json:"totalExercises"`
	TotalRoutines  int64 `json:"totalRoutines"`
	TotalWorkouts  int64 `json:"totalWorkouts"`
}

type AdminStatsResponse struct {
	Totals              AdminTotals              `json:"totals"`
	UsersByMonth        []WeeklyCount            `json:"usersByMonth"`
	ExercisesByCategory []RoutineCount           `json:"exercisesByCategory"`
	PopularRoutines     []RoutineCount           `json:"popularRoutines"`
	RecentActivity      []map[string]interface{} `json:"recentActivity"`
}
