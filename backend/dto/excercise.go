package dto

type ExerciseRequest struct {
	UserID      string   `json:"user_id"`
	Name        string   `json:"name" binding:"required"`
	Category    string   `json:"category" binding:"required"`
	MuscleGroup string   `json:"muscle_group" binding:"required"`
	Description string   `json:"description,omitempty"`
	Difficulty  string   `json:"difficulty" binding:"required"`
	MediaURL    string   `json:"media_url,omitempty"`
	Steps       []string `json:"steps,omitempty"`
}

type ExerciseResponse struct {
	ID          string   `json:"id"`
	UserID      string   `json:"user_id"`
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description,omitempty"`
	Category    string   `json:"category" binding:"required"`
	MuscleGroup string   `json:"muscle_group" binding:"required"`
	Difficulty  string   `json:"difficulty" binding:"required"`
	MediaURL    string   `json:"media_url,omitempty"`
	Steps       []string `json:"steps,omitempty"`
}

type ExerciseListResponse struct {
	Exercises []ExerciseResponse `json:"exercises"`
}

type ExerciseSearch struct {
	Name        string `form:"name"`
	Category    string `form:"category"`
	MuscleGroup string `form:"muscle_group"`
}
