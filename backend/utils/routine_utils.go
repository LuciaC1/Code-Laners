package utils

import (
	"backend/dto"
	"backend/models"
)

func ConverModelToRoutineDTO(routine models.Routine) dto.RoutineDTO {
	entries := make([]dto.RoutineEntryDTO, len(routine.Entries))
	for i, entry := range routine.Entries {
		entries[i] = dto.RoutineEntryDTO{
			ExerciseID: entry.ExerciseID.Hex(),
			Order:      entry.Order,
			Sets:       entry.Sets,
			Reps:       entry.Reps,
			Weight:     entry.Weight,
		}
	}
	return dto.RoutineDTO{
		Name:        routine.Name,
		Description: routine.Description,
		Entries:     entries,
		IsPublic:    routine.IsPublic,
	}
}
func ConvertModelToRoutineEntryDTO(entry models.RoutineEntry) dto.RoutineEntryDTO {
	return dto.RoutineEntryDTO{
		ExerciseID: entry.ExerciseID.Hex(),
		Order:      entry.Order,
		Sets:       entry.Sets,
		Reps:       entry.Reps,
		Weight:     entry.Weight,
	}
}
