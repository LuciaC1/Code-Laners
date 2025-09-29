package dto

type RoutineEntryDTO struct {
	ExerciseID string  `json:"exercise_id" binding:"required"`
	Order      int     `json:"order" binding:"required"`
	Sets       int     `json:"sets" binding:"required"`
	Reps       int     `json:"reps" binding:"required"`
	Weight     float64 `json:"weight,omitempty"`
}

type RoutineDTO struct {
	Name        string            `json:"name" binding:"required"`
	Description string            `json:"description,omitempty"`
	Entries     []RoutineEntryDTO `json:"entries" binding:"required,dive,required"`
	IsPublic    bool              `json:"is_public"`
}
