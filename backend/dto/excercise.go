package dto

type ExerciseDTO struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description,omitempty"`
	Category    string   `json:"category" binding:"required"`
	MuscleGroup string   `json:"muscle_group" binding:"required"`
	Difficulty  string   `json:"difficulty" binding:"required"`
	MediaURL    string   `json:"media_url,omitempty"`
	Steps       []string `json:"steps,omitempty"`
}

type ExerciseQuery struct {
	Name     string `form:"name"`
	Category string `form:"category"`
}
