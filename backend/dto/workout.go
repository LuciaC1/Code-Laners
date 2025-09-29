package dto

type PerformedExerciseDTO struct {
	ExerciseID string  `json:"exercise_id" binding:"required"`
	Sets       int     `json:"sets" binding:"required"`
	Reps       int     `json:"reps" binding:"required"`
	Weight     float64 `json:"weight,omitempty"`
}

type WorkoutDTO struct {
	RoutineID   string                 `json:"routine_id,omitempty"`
	Performed   []PerformedExerciseDTO `json:"performed" binding:"required"`
	DurationMin int                    `json:"duration_min,omitempty"`
	Notes       string                 `json:"notes,omitempty"`
}
