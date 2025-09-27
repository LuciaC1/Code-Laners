package services

import (
	"errors"
	"time"

	"backend/models"
	"backend/repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

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

func (s *ExerciseService) CreateExercise(e models.Exercise, creatorIDHex, actorRole string) (*mongo.InsertOneResult, error) {
	if actorRole != "admin" {
		return nil, errors.New("forbidden: only admins can create exercises")
	}
	if e.Name == "" {
		return nil, errors.New("exercise name is required")
	}
	if e.Category == "" {
		return nil, errors.New("category is required")
	}
	if e.MuscleGroup == "" {
		return nil, errors.New("muscle group is required")
	}
	if e.Difficulty == "" {
		return nil, errors.New("difficulty is required")
	}

	creatorOID, err := primitive.ObjectIDFromHex(creatorIDHex)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	e.CreatedBy = creatorOID
	e.CreatedAt = now
	e.UpdatedAt = now

	return s.repo.CreateExercise(e)
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
