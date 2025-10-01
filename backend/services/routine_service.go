package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"backend/dto"
	"backend/models"
	"backend/repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoutineResponseDTO struct {
	ID          string                `json:"id"`
	OwnerID     string                `json:"owner_id"`
	Name        string                `json:"name"`
	Description string                `json:"description,omitempty"`
	Entries     []dto.RoutineEntryDTO `json:"entries"`
	IsPublic    bool                  `json:"is_public"`
	CreatedAt   time.Time             `json:"created_at"`
	UpdatedAt   time.Time             `json:"updated_at"`
}

type RoutineServiceInterface interface {
	CreateRoutine(ownerID string, input dto.RoutineDTO) (string, error)
	GetRoutines(ownerID string, name string) ([]RoutineResponseDTO, error)
	GetRoutineByID(id string) (RoutineResponseDTO, error)
	UpdateRoutine(ownerID string, routineID string, input dto.RoutineDTO) error
	DeleteRoutine(ownerID string, routineID string) error
	DuplicateRoutine(ownerID string, sourceRoutineID string, newName string) (string, error)
}

type RoutineService struct {
	repo         repositories.RoutineRepositoryInterface
	exerciseRepo repositories.ExerciseRepositoryInterface
}

func NewRoutineService(repo repositories.RoutineRepositoryInterface, exerciseRepo repositories.ExerciseRepositoryInterface) *RoutineService {
	return &RoutineService{repo: repo, exerciseRepo: exerciseRepo}
}

func (s *RoutineService) CreateRoutine(ownerID string, input dto.RoutineDTO) (string, error) {
	if ownerID == "" {
		return "", errors.New("ownerID requerido")
	}
	own, err := primitive.ObjectIDFromHex(ownerID)
	if err != nil {
		return "", fmt.Errorf("ownerID inválido: %w", err)
	}
	if err := validateRoutineEntries(input.Entries); err != nil {
		return "", err
	}
	var exIDs []primitive.ObjectID
	for _, e := range input.Entries {
		id, _ := primitive.ObjectIDFromHex(e.ExerciseID)
		exIDs = append(exIDs, id)
	}
	if err := s.verifyExercisesExist(exIDs); err != nil {
		return "", err
	}
	routine := models.Routine{
		ID:          primitive.NewObjectID(),
		OwnerID:     own,
		Name:        input.Name,
		Description: input.Description,
		IsPublic:    input.IsPublic,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	routine.Entries = make([]models.RoutineEntry, 0, len(input.Entries))
	for _, e := range input.Entries {
		exID, _ := primitive.ObjectIDFromHex(e.ExerciseID)
		routine.Entries = append(routine.Entries, models.RoutineEntry{
			ExerciseID: exID,
			Order:      e.Order,
			Sets:       e.Sets,
			Reps:       e.Reps,
			Weight:     e.Weight,
		})
	}

	res, err := s.repo.CreateRoutine(routine)
	if err != nil {
		return "", err
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		return oid.Hex(), nil
	}
	return "", nil
}

func (s *RoutineService) GetRoutines(ownerID string, name string) ([]RoutineResponseDTO, error) {
	if ownerID == "" {
		return nil, errors.New("ownerID requerido")
	}
	own, err := primitive.ObjectIDFromHex(ownerID)
	if err != nil {
		return nil, err
	}
	modelsList, err := s.repo.GetRoutines(own, name)
	if err != nil {
		return nil, err
	}
	out := make([]RoutineResponseDTO, 0, len(modelsList))
	for _, m := range modelsList {
		out = append(out, routineModelToResponse(m))
	}
	return out, nil
}

func (s *RoutineService) GetRoutineByID(id string) (RoutineResponseDTO, error) {
	m, err := s.repo.GetRoutineByID(id)
	if err != nil {
		return RoutineResponseDTO{}, err
	}
	return routineModelToResponse(m), nil
}

func (s *RoutineService) UpdateRoutine(ownerID string, routineID string, input dto.RoutineDTO) error {
	existing, err := s.repo.GetRoutineByID(routineID)
	if err != nil {
		return err
	}
	own, err := primitive.ObjectIDFromHex(ownerID)
	if err != nil {
		return err
	}
	if existing.OwnerID != own {
		return errors.New("no autorizado: no es el owner de la rutina")
	}
	if err := validateRoutineEntries(input.Entries); err != nil {
		return err
	}
	var exIDs []primitive.ObjectID
	for _, e := range input.Entries {
		id, _ := primitive.ObjectIDFromHex(e.ExerciseID)
		exIDs = append(exIDs, id)
	}
	if err := s.verifyExercisesExist(exIDs); err != nil {
		return err
	}
	existing.Name = input.Name
	existing.Description = input.Description
	existing.IsPublic = input.IsPublic
	existing.UpdatedAt = time.Now()
	existing.Entries = make([]models.RoutineEntry, 0, len(input.Entries))
	for _, e := range input.Entries {
		exID, _ := primitive.ObjectIDFromHex(e.ExerciseID)
		existing.Entries = append(existing.Entries, models.RoutineEntry{
			ExerciseID: exID,
			Order:      e.Order,
			Sets:       e.Sets,
			Reps:       e.Reps,
			Weight:     e.Weight,
		})
	}
	_, err = s.repo.UpdateRoutine(existing)
	return err
}

func (s *RoutineService) DeleteRoutine(ownerID string, routineID string) error {
	existing, err := s.repo.GetRoutineByID(routineID)
	if err != nil {
		return err
	}
	own, err := primitive.ObjectIDFromHex(ownerID)
	if err != nil {
		return err
	}
	if existing.OwnerID != own {
		return errors.New("no autorizado: no es el owner de la rutina")
	}
	_, err = s.repo.DeleteRoutine(existing.ID)
	return err
}

func (s *RoutineService) DuplicateRoutine(ownerID string, sourceRoutineID string, newName string) (string, error) {
	if sourceRoutineID == "" {
		return "", errors.New("sourceRoutineID requerido")
	}
	src, err := s.repo.GetRoutineByID(sourceRoutineID)
	if err != nil {
		return "", err
	}
	own, err := primitive.ObjectIDFromHex(ownerID)
	if err != nil {
		return "", fmt.Errorf("ownerID inválido: %w", err)
	}
	var exIDs []primitive.ObjectID
	for _, e := range src.Entries {
		exIDs = append(exIDs, e.ExerciseID)
	}
	if err := s.verifyExercisesExist(exIDs); err != nil {
		return "", err
	}
	copy := models.Routine{
		ID:          primitive.NewObjectID(),
		OwnerID:     own,
		Name:        newName,
		Description: src.Description,
		IsPublic:    src.IsPublic,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	copy.Entries = make([]models.RoutineEntry, 0, len(src.Entries))
	for _, e := range src.Entries {
		copy.Entries = append(copy.Entries, models.RoutineEntry{
			ExerciseID: e.ExerciseID,
			Order:      e.Order,
			Sets:       e.Sets,
			Reps:       e.Reps,
			Weight:     e.Weight,
		})
	}
	res, err := s.repo.CreateRoutine(copy)
	if err != nil {
		return "", err
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		return oid.Hex(), nil
	}
	return "", nil
}

func validateRoutineEntries(entries []dto.RoutineEntryDTO) error {
	if len(entries) == 0 {
		return errors.New("la rutina debe contener al menos un ejercicio")
	}
	orders := make(map[int]bool)
	for i, e := range entries {
		if e.ExerciseID == "" {
			return fmt.Errorf("entry %d: exercise_id requerido", i)
		}
		if _, err := primitive.ObjectIDFromHex(e.ExerciseID); err != nil {
			return fmt.Errorf("entry %d: exercise_id inválido: %w", i, err)
		}
		if e.Order <= 0 {
			return fmt.Errorf("entry %d: order debe ser > 0", i)
		}
		if orders[e.Order] {
			return fmt.Errorf("entry %d: order duplicado (%d)", i, e.Order)
		}
		orders[e.Order] = true
		if e.Sets <= 0 {
			return fmt.Errorf("entry %d: sets debe ser > 0", i)
		}
		if e.Reps <= 0 {
			return fmt.Errorf("entry %d: reps debe ser > 0", i)
		}
	}
	return nil
}

func (s *RoutineService) verifyExercisesExist(ids []primitive.ObjectID) error {
	if len(ids) == 0 {
		return nil
	}
	missing := make([]string, 0)
	seen := make(map[string]bool)
	for _, id := range ids {
		if id.IsZero() {
			continue
		}
		h := id.Hex()
		if seen[h] {
			continue
		}
		seen[h] = true
		exists, err := s.exerciseRepo.GetExerciseByID(h)
		if err != nil {
			return err
		}
		if exists.ID.IsZero() {
			missing = append(missing, h)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("exercises not found: %s", strings.Join(missing, ","))
	}
	return nil
}

func routineModelToResponse(m models.Routine) RoutineResponseDTO {
	entries := make([]dto.RoutineEntryDTO, 0, len(m.Entries))
	for _, e := range m.Entries {
		entries = append(entries, dto.RoutineEntryDTO{
			ExerciseID: e.ExerciseID.Hex(),
			Order:      e.Order,
			Sets:       e.Sets,
			Reps:       e.Reps,
			Weight:     e.Weight,
		})
	}
	return RoutineResponseDTO{
		ID:          m.ID.Hex(),
		OwnerID:     m.OwnerID.Hex(),
		Name:        m.Name,
		Description: m.Description,
		Entries:     entries,
		IsPublic:    m.IsPublic,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}
