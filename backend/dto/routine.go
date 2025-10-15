package dto

type RoutineRequest struct {
	UserID      string                 `json:"user_id"`
	Name        string                 `json:"name" binding:"required"`
	Excercises  []RoutineExcerciseList `json:"exercises" binding:"required"`
	Description string                 `json:"description,omitempty"`
	IsPublic    bool                   `json:"is_public"`
}

type RoutineResponse struct {
	ID          string                 `json:"id"`
	UserID      string                 `json:"user_id"`
	Name        string                 `json:"name" binding:"required"`
	Excercises  []RoutineExcerciseList `json:"exercises" binding:"required"`
	Description string                 `json:"description,omitempty"`
	IsPublic    bool                   `json:"is_public"`
}

type RoutineExcerciseList struct {
	ExerciseID string  `json:"exercise_id" binding:"required"`
	Order      int     `json:"order" binding:"required"`
	Sets       int     `json:"sets" binding:"required"`
	Reps       int     `json:"reps" binding:"required"`
	Weight     float64 `json:"weight,omitempty"`
}
