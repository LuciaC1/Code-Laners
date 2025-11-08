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
type RoutineServiceInterface interface {
	CreateRoutine(ownerID string, input dto.RoutineRequest) (dto.RoutineResponse, error)
	GetRoutines(ownerID string, name string) ([]dto.RoutineResponse, error)
	GetRoutineByID(id string) (dto.RoutineResponse, error)
	UpdateRoutine(ownerID string, routineID string, input dto.RoutineRequest) (dto.RoutineResponse, error)
	DeleteRoutine(ownerID string, routineID string) error
	DuplicateRoutine(ownerID string, sourceRoutineID string, newName string) (string, error)
}
type RoutineService struct {
	repo         repositories.RoutineRepositoryInterface
	exerciseRepo repositories.ExerciseRepositoryInterface
	userRepo     repositories.UserRepositoryInterface
}
func NewRoutineService(repo repositories.RoutineRepositoryInterface, exerciseRepo repositories.ExerciseRepositoryInterface, userRepo repositories.UserRepositoryInterface) *RoutineService {
	return &RoutineService{
		repo:         repo,
		exerciseRepo: exerciseRepo,
		userRepo:     userRepo,
	}
}
func (s *RoutineService) CreateRoutine(ownerID string, input dto.RoutineRequest) (dto.RoutineResponse, error) {
	if ownerID == "" {
		return dto.RoutineResponse{}, errors.New("ownerID requerido")
	}
	own, err := primitive.ObjectIDFromHex(ownerID)
	if err != nil {
		return dto.RoutineResponse{}, fmt.Errorf("ownerID inválido: %w", err)
	}
	if err := validateRoutineEntries(input.Excercises); err != nil {
		return dto.RoutineResponse{}, err
	}
	var exIDs []primitive.ObjectID
	for _, e := range input.Excercises {
		id, _ := primitive.ObjectIDFromHex(e.ExerciseID)
		exIDs = append(exIDs, id)
	}
	if err := s.verifyExercisesExist(exIDs); err != nil {
		return dto.RoutineResponse{}, err
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
	routine.Entries = make([]models.RoutineExcerciseList, 0, len(input.Excercises))
	for _, e := range input.Excercises {
		exID, _ := primitive.ObjectIDFromHex(e.ExerciseID)
		routine.Entries = append(routine.Entries, models.RoutineExcerciseList{
			ExerciseID: exID,
			Order:      e.Order,
			Sets:       e.Sets,
			Reps:       e.Reps,
			Weight:     e.Weight,
		})
	}
	res, err := s.repo.CreateRoutine(routine)
	if err != nil {
		return dto.RoutineResponse{}, err
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		routine.ID = oid
		return s.convertModelToDTOWithExerciseNames(routine), nil
	}
	return dto.RoutineResponse{}, nil
}
func (s *RoutineService) GetRoutines(ownerID string, name string) ([]dto.RoutineResponse, error) {
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
	out := make([]dto.RoutineResponse, 0, len(modelsList))
	for _, m := range modelsList {
		out = append(out, s.convertModelToDTOWithExerciseNames(m))
	}
	return out, nil
}
func (s *RoutineService) GetRoutineByID(id string) (dto.RoutineResponse, error) {
	m, err := s.repo.GetRoutineByID(id)
	if err != nil {
		return dto.RoutineResponse{}, err
	}
	return s.convertModelToDTOWithExerciseNames(m), nil
}
func (s *RoutineService) UpdateRoutine(ownerID string, routineID string, input dto.RoutineRequest) (dto.RoutineResponse, error) {
	existing, err := s.repo.GetRoutineByID(routineID)
	if err != nil {
		return dto.RoutineResponse{}, err
	}
	own, err := primitive.ObjectIDFromHex(ownerID)
	if err != nil {
		return dto.RoutineResponse{}, err
	}
	if existing.OwnerID != own {
		return dto.RoutineResponse{}, errors.New("no autorizado: no es el owner de la rutina")
	}
	if err := validateRoutineEntries(input.Excercises); err != nil {
		return dto.RoutineResponse{}, err
	}
	var exIDs []primitive.ObjectID
	for _, e := range input.Excercises {
		id, _ := primitive.ObjectIDFromHex(e.ExerciseID)
		exIDs = append(exIDs, id)
	}
	if err := s.verifyExercisesExist(exIDs); err != nil {
		return dto.RoutineResponse{}, err
	}
	existing.Name = input.Name
	existing.Description = input.Description
	existing.IsPublic = input.IsPublic
	existing.UpdatedAt = time.Now()
	existing.Entries = make([]models.RoutineExcerciseList, 0, len(input.Excercises))
	for _, e := range input.Excercises {
		exID, _ := primitive.ObjectIDFromHex(e.ExerciseID)
		existing.Entries = append(existing.Entries, models.RoutineExcerciseList{
			ExerciseID: exID,
			Order:      e.Order,
			Sets:       e.Sets,
			Reps:       e.Reps,
			Weight:     e.Weight,
		})
	}
	_, err = s.repo.UpdateRoutine(existing)
	if err != nil {
		return dto.RoutineResponse{}, err
	}
	updated, err := s.repo.GetRoutineByID(routineID)
	if err != nil {
		return dto.RoutineResponse{}, err
	}
	return s.convertModelToDTOWithExerciseNames(updated), nil
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
	copy.Entries = make([]models.RoutineExcerciseList, 0, len(src.Entries))
	for _, e := range src.Entries {
		copy.Entries = append(copy.Entries, models.RoutineExcerciseList{
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
func validateRoutineEntries(entries []dto.RoutineExcerciseList) error {
	orders := make(map[int]bool)
	for i, e := range entries {
		if _, err := primitive.ObjectIDFromHex(e.ExerciseID); err != nil {
			return fmt.Errorf("entry %d: exercise_id inválido: %w", i, err)
		}
		if orders[e.Order] {
			return fmt.Errorf("entry %d: order duplicado (%d)", i, e.Order)
		}
		orders[e.Order] = true
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
func (s *RoutineService) convertModelToDTOWithExerciseNames(routine models.Routine) dto.RoutineResponse {
	entries := make([]dto.RoutineExcerciseList, len(routine.Entries))
	for i, entry := range routine.Entries {
		exerciseName := ""
		if !entry.ExerciseID.IsZero() {
			exercise, err := s.exerciseRepo.GetExerciseByID(entry.ExerciseID.Hex())
			if err == nil && !exercise.ID.IsZero() {
				exerciseName = exercise.Name
			}
		}
		entries[i] = dto.RoutineExcerciseList{
			ExerciseID:   entry.ExerciseID.Hex(),
			ExerciseName: exerciseName,
			Order:        entry.Order,
			Sets:         entry.Sets,
			Reps:         entry.Reps,
			Weight:       entry.Weight,
		}
	}
	ownerName := ""
	if s.userRepo != nil && !routine.OwnerID.IsZero() {
		user, err := s.userRepo.GetUserByID(routine.OwnerID.Hex())
		if err == nil && !user.ID.IsZero() {
			ownerName = user.Name
		}
	}
	return dto.RoutineResponse{
		ID:          routine.ID.Hex(),
		UserID:      routine.OwnerID.Hex(),
		OwnerName:   ownerName,
		Name:        routine.Name,
		Excercises:  entries,
		Description: routine.Description,
		IsPublic:    routine.IsPublic,
	}
}
