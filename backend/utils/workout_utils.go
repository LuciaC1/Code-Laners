package utils

import (
	"backend/dto"
	"backend/models"
)

func ConvertWorkoutModelToDTO(workout models.Workout) dto.WorkoutDTO {
	return dto.WorkoutDTO{
		ID:                workout.ID,
		UserID:            workout.UserID.Hex(),
		RoutineID:         workout.RoutineID.Hex(),
		CompletedAt:       workout.CompletedAt,
		UpdatedAt:         workout.UpdatedAt,
		DurationMinutes:   workout.DurationMinutes,
		Notes:             workout.Notes,
		EstimatedCalories: workout.EstimatedCalories,
	}
}
