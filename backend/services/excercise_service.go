package services

import (
	"errors"
	"time"

	"backend/dto"
	"backend/models"
	"backend/repositories"
	"backend/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ExerciseInterface interface {
	GetExercises(name, category, muscleGroup string) ([]models.Exercise, error)
	GetExerciseByID(id string) (models.Exercise, error)
	CreateExercise(exercise dto.ExerciseRequest) (dto.ExerciseResponse, error)
	UpdateExercise(id string, exercise dto.ExerciseRequest) (dto.ExerciseResponse, error)
	DeleteExercise(idHex, actorRole string) (*mongo.DeleteResult, error)
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
func (service *ExerciseService) CreateExercise(exercise dto.ExerciseRequest) (dto.ExerciseResponse, error) {
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

func (s *ExerciseService) DeleteExercise(idHex, actorRole string) (*mongo.DeleteResult, error) {
	if actorRole != "admin" {
		return nil, errors.New("forbidden: only admins can delete exercises")
	}
	if idHex == "" {
		return nil, errors.New("id required")
	}
	deletedId, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return nil, err
	}
	return s.repo.DeleteExercise(deletedId)
}
