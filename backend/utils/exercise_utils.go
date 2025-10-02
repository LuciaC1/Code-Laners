package utils

import (
	"backend/dto"
	"backend/models"
)

func ConvertExerciseModelToDTO(exercise models.Exercise) dto.ExerciseDTO {
	return dto.ExerciseDTO{
		UserID:      exercise.UserID,
		Name:        exercise.Name,
		Description: exercise.Description,
		Category:    exercise.Category,
		MuscleGroup: exercise.MuscleGroup,
		Difficulty:  exercise.Difficulty,
		MediaURL:    exercise.MediaURL,
		Steps:       exercise.Steps,
	}
}
func ConvertExerciseModelsToDTOs(exercises []models.Exercise) []dto.ExerciseDTO {
	dtos := make([]dto.ExerciseDTO, len(exercises))
	for i, exercise := range exercises {
		dtos[i] = ConvertExerciseModelToDTO(exercise)
	}
	return dtos
}
