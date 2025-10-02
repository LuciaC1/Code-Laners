package dto

type ExerciseDTO struct {
	UserID      string   `json:"user_id"`
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description,omitempty"`
	Category    string   `json:"category" binding:"required"`
	MuscleGroup string   `json:"muscle_group" binding:"required"`
	Difficulty  string   `json:"difficulty" binding:"required"`
	MediaURL    string   `json:"media_url,omitempty"`
	Steps       []string `json:"steps,omitempty"`
}

type ExerciseRequest struct {
	Name        string `form:"name"`
	Category    string `form:"category"`
	MuscleGroup string `form:"muscle_group"`
}

type ExerciseListResponse struct {
	Exercises []ExerciseDTO `json:"exercises"`
}
