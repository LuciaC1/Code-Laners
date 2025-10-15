package services

import (
	"errors"
	"time"

	"backend/dto"
	"backend/models"
	"backend/repositories"
	"backend/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ExerciseInterface interface {
	GetExercises(name, category, muscleGroup string) ([]models.Exercise, error)
	GetExerciseByID(id string) (models.Exercise, error)
	CreateExercise(exercise dto.ExerciseRequest) (dto.ExerciseResponse, error)
	UpdateExercise(id string, exercise dto.ExerciseRequest) (dto.ExerciseResponse, error)
	DeleteExercise(ownerID, exerciseID string) error
	SearchExercises(search dto.ExerciseSearch) ([]dto.ExerciseResponse, error) // new function signature
}

type ExerciseService struct {
	repo repositories.ExerciseRepositoryInterface
}

func NewExerciseService(repo repositories.ExerciseRepositoryInterface) *ExerciseService {
	return &ExerciseService{repo: repo}
}

func (s *ExerciseService) GetExercises(name, category, muscleGroup string) ([]models.Exercise, error) {
	return s.repo.GetExercises(name, category, muscleGroup)
}

func (s *ExerciseService) GetExerciseByID(id string) (models.Exercise, error) {
	if id == "" {
		return models.Exercise{}, errors.New("id required")
	}
	return s.repo.GetExerciseByID(id)
}

func validateExerciseRequest(request dto.ExerciseRequest) error {
	if request.UserID == "" {
		return errors.New("user id is required")
	}
	if request.Name == "" {
		return errors.New("name is required")
	}
	if request.Category == "" {
		return errors.New("category is required")
	}
	if request.MuscleGroup == "" {
		return errors.New("muscle group is required")
	}
	if request.Difficulty == "" {
		return errors.New("difficulty is required")
	}

	if _, err := primitive.ObjectIDFromHex(request.UserID); err != nil {
		return errors.New("invalid user id")
	}
	return nil
}

func (service *ExerciseService) CreateExercise(exercise dto.ExerciseRequest) (dto.ExerciseResponse, error) {

	if err := validateExerciseRequest(exercise); err != nil {
		return dto.ExerciseResponse{}, err
	}

	modelExercise := utils.ConvertRequestToExerciseModel(exercise)

	result, err := service.repo.CreateExercise(modelExercise)
	if err != nil {
		return dto.ExerciseResponse{}, err
	}
	createdExercise, err := service.repo.GetExerciseByID(result.InsertedID.(primitive.ObjectID).Hex())
	if err != nil {
		return dto.ExerciseResponse{}, err
	}
	return utils.ConvertExerciseModelToDTO(createdExercise), nil
}

func (service *ExerciseService) UpdateExercise(id string, exercise dto.ExerciseRequest) (dto.ExerciseResponse, error) {
	if err := validateExerciseRequest(exercise); err != nil {
		return dto.ExerciseResponse{}, err
	}

	existing, err := service.repo.GetExerciseByID(id)
	if err != nil {
		return dto.ExerciseResponse{}, err
	}
	if existing.UserID != exercise.UserID {
		return dto.ExerciseResponse{}, errors.New("unauthorized: cannot update exercise you do not own")
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return dto.ExerciseResponse{}, errors.New("invalid id")
	}

	modelExercise := utils.ConvertRequestToExerciseModel(exercise)
	modelExercise.ID = objectID
	modelExercise.UpdatedAt = time.Now()

	_, err = service.repo.UpdateExercise(modelExercise)
	if err != nil {
		return dto.ExerciseResponse{}, err
	}

	return utils.ConvertExerciseModelToDTO(modelExercise), nil
}

func (s *ExerciseService) DeleteExercise(ownerID, exerciseID string) error {
	if ownerID == "" {
		return errors.New("owner id is required")
	}
	existing, err := s.repo.GetExerciseByID(exerciseID)
	if err != nil {
		return err
	}
	if existing.UserID != ownerID {
		return errors.New("unauthorized: cannot delete exercise you do not own")
	}
	objID, err := primitive.ObjectIDFromHex(exerciseID)
	if err != nil {
		return err
	}
	_, err = s.repo.DeleteExercise(objID)
	return err
}

func (s *ExerciseService) SearchExercises(search dto.ExerciseSearch) ([]dto.ExerciseResponse, error) {
	exercises, err := s.repo.GetExercises(search.Name, search.Category, search.MuscleGroup)
	if err != nil {
		return nil, err
	}
	return utils.ConvertExerciseModelsToDTOList(exercises), nil
}
