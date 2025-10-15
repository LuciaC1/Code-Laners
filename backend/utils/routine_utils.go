package utils

import (
	"backend/dto"
	"backend/models"
)

func ConverModelToRoutineDTO(routine models.Routine) dto.RoutineResponse {
	entries := make([]dto.RoutineExcerciseList, len(routine.Entries))
	for i, entry := range routine.Entries {
		entries[i] = dto.RoutineExcerciseList{
			ExerciseID: entry.ExerciseID.Hex(),
			Order:      entry.Order,
			Sets:       entry.Sets,
			Reps:       entry.Reps,
			Weight:     entry.Weight,
		}
	}
	return dto.RoutineResponse{
		ID:          routine.ID.Hex(),
		UserID:      routine.OwnerID.Hex(),
		Name:        routine.Name,
		Excercises:  entries,
		Description: routine.Description,
		IsPublic:    routine.IsPublic,
	}
}

func ConvertModelToRoutineExcerciseListDTO(entry models.RoutineExcerciseList) dto.RoutineExcerciseList {
	return dto.RoutineExcerciseList{
		ExerciseID: entry.ExerciseID.Hex(),
		Order:      entry.Order,
		Sets:       entry.Sets,
		Reps:       entry.Reps,
		Weight:     entry.Weight,
	}
}
