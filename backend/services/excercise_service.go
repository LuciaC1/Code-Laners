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
	CreateExercise(e models.Exercise, creatorIDHex, actorRole string) (*mongo.InsertOneResult, error)
	UpdateExercise(idHex string, payload models.Exercise, actorRole string) (*mongo.UpdateResult, error)
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
func (service *ExerciseService) CreateExercise(exercise dto.ExerciseDTO) (dto.ExerciseDTO, error) {
	modelExercise := utils.ConvertRequestToExerciseModel(exercise)
	result, err := service.repo.CreateExercise(modelExercise)
	if err != nil {
		return dto.ExerciseDTO{}, err
	}
	createdExercise, err := service.repo.GetExerciseByID(result.InsertedID.(primitive.ObjectID).Hex())
	if err != nil {
		return dto.ExerciseDTO{}, err
	}
	return utils.ConvertExerciseModelToDTO(createdExercise), nil
}

func (s *ExerciseService) UpdateExercise(idHex string, payload models.Exercise, actorRole string) (*mongo.UpdateResult, error) {
	if actorRole != "admin" {
		return nil, errors.New("forbidden: only admins can update exercises")
	}
	if idHex == "" {
		return nil, errors.New("id required")
	}

	existing, err := s.repo.GetExerciseByID(idHex)
	if err != nil {
		return nil, err
	}

	if payload.Name != "" {
		existing.Name = payload.Name
	}
	if payload.Description != "" {
		existing.Description = payload.Description
	}
	if payload.Category != "" {
		existing.Category = payload.Category
	}
	if payload.MuscleGroup != "" {
		existing.MuscleGroup = payload.MuscleGroup
	}
	if payload.Difficulty != "" {
		existing.Difficulty = payload.Difficulty
	}
	if payload.MediaURL != "" {
		existing.MediaURL = payload.MediaURL
	}
	if len(payload.Steps) > 0 {
		existing.Steps = payload.Steps
	}
	existing.UpdatedAt = time.Now()

	return s.repo.UpdateExercise(existing)
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
