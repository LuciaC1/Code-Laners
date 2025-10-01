package services

import (
	"errors"
	"time"

	"backend/dto"
	"backend/models"
	"backend/repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WorkoutServiceInterface interface {
	GetWorkouts(userID string) ([]dto.WorkoutDTO, error)
	GetWorkoutByID(id string) (dto.WorkoutDTO, error)
	CreateWorkout(input dto.WorkoutDTO) (string, error)
	UpdateWorkout(input dto.WorkoutDTO) error
	DeleteWorkout(id string) error
}

type WorkoutService struct {
	repo repositories.WorkoutRepositoryInterface
}

func NewWorkoutService(repo repositories.WorkoutRepositoryInterface) *WorkoutService {
	return &WorkoutService{repo: repo}
}

func (s *WorkoutService) GetWorkouts(userID string) ([]dto.WorkoutDTO, error) {
	if userID == "" {
		return nil, errors.New("userID requerido")
	}
	uid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}
	modelsList, err := s.repo.GetWorkouts(uid)
	if err != nil {
		return nil, err
	}
	var dtos []dto.WorkoutDTO
	for _, m := range modelsList {
		dtos = append(dtos, modelToDTO(m))
	}
	return dtos, nil
}

func (s *WorkoutService) GetWorkoutByID(id string) (dto.WorkoutDTO, error) {
	m, err := s.repo.GetWorkoutByID(id)
	if err != nil {
		return dto.WorkoutDTO{}, err
	}
	return modelToDTO(m), nil
}

func (s *WorkoutService) CreateWorkout(input dto.WorkoutDTO) (string, error) {
	if input.UserID == "" {
		return "", errors.New("user_id requerido")
	}
	uid, err := primitive.ObjectIDFromHex(input.UserID)
	if err != nil {
		return "", err
	}

	var rid primitive.ObjectID
	if input.RoutineID != "" {
		rid, err = primitive.ObjectIDFromHex(input.RoutineID)
		if err != nil {
			return "", err
		}
	}

	workout := models.Workout{
		ID:                primitive.NewObjectID(),
		UserID:            uid,
		RoutineID:         rid,
		CompletedAt:       input.CompletedAt,
		UpdatedAt:         time.Now(),
		DurationMinutes:   input.DurationMinutes,
		EstimatedCalories: input.EstimatedCalories,
		Notes:             input.Notes,
	}

	res, err := s.repo.CreateWorkout(workout)
	if err != nil {
		return "", err
	}
	if res == nil {
		return "", errors.New("insert result nil")
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		return oid.Hex(), nil
	}
	return "", nil
}

func (s *WorkoutService) UpdateWorkout(input dto.WorkoutDTO) error {
	if input.ID.IsZero() {
		return errors.New("id requerido para actualizar")
	}
	var uid primitive.ObjectID
	var err error
	if input.UserID != "" {
		uid, err = primitive.ObjectIDFromHex(input.UserID)
		if err != nil {
			return err
		}
	}
	var rid primitive.ObjectID
	if input.RoutineID != "" {
		rid, err = primitive.ObjectIDFromHex(input.RoutineID)
		if err != nil {
			return err
		}
	}

	workout := models.Workout{
		ID:                input.ID,
		UserID:            uid,
		RoutineID:         rid,
		CompletedAt:       input.CompletedAt,
		UpdatedAt:         time.Now(),
		DurationMinutes:   input.DurationMinutes,
		EstimatedCalories: input.EstimatedCalories,
		Notes:             input.Notes,
	}

	_, err = s.repo.UpdateWorkout(workout)
	return err
}

func (s *WorkoutService) DeleteWorkout(id string) error {
	if id == "" {
		return errors.New("id requerido")
	}
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = s.repo.DeleteWorkout(objID)
	return err
}

func modelToDTO(m models.Workout) dto.WorkoutDTO {
	var uidHex string
	if !m.UserID.IsZero() {
		uidHex = m.UserID.Hex()
	}
	var ridHex string
	if !m.RoutineID.IsZero() {
		ridHex = m.RoutineID.Hex()
	}
	return dto.WorkoutDTO{
		ID:                m.ID,
		UserID:            uidHex,
		RoutineID:         ridHex,
		CompletedAt:       m.CompletedAt,
		UpdatedAt:         m.UpdatedAt,
		DurationMinutes:   m.DurationMinutes,
		Notes:             m.Notes,
		EstimatedCalories: m.EstimatedCalories,
	}
}

