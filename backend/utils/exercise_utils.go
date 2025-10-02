package utils

import (
	"backend/dto"
	"backend/models"
	"time"
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
func ConvertExerciseModelsToDTOList(exercises []models.Exercise) []dto.ExerciseDTO {
	dtos := make([]dto.ExerciseDTO, len(exercises))
	for i, exercise := range exercises {
		dtos[i] = ConvertExerciseModelToDTO(exercise)
	}
	return dtos
}

func ConvertModelToExerciseRequest(exercise models.Exercise) dto.ExerciseRequest {
	return dto.ExerciseRequest{
		Name:        exercise.Name,
		Category:    exercise.Category,
		MuscleGroup: exercise.MuscleGroup,
	}
}
func ConvertRequestToExerciseModel(request dto.ExerciseDTO) models.Exercise {
	return models.Exercise{
		UserID:      request.UserID,
		Name:        request.Name,
		Description: request.Description,
		Category:    request.Category,
		MuscleGroup: request.MuscleGroup,
		Difficulty:  request.Difficulty,
		MediaURL:    request.MediaURL,
		Steps:       request.Steps,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}
