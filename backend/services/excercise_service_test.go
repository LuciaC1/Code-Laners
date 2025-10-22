package services

import (
	"testing"
	"time"

	"backend/dto"
	"backend/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mockRepo struct {
	getExercisesFn func(name, category, muscleGroup string) ([]models.Exercise, error)
	getByIDFn      func(id string) (models.Exercise, error)
	createFn       func(ex models.Exercise) (primitive.ObjectID, error)
	updateFn       func(ex models.Exercise) (bool, error)
	deleteFn       func(id primitive.ObjectID) (bool, error)
}

func (m *mockRepo) GetExercises(name, category, muscleGroup string) ([]models.Exercise, error) {
	return m.getExercisesFn(name, category, muscleGroup)
}
func (m *mockRepo) GetExerciseByID(id string) (models.Exercise, error) {
	return m.getByIDFn(id)
}
func (m *mockRepo) CreateExercise(ex models.Exercise) (*mongo.InsertOneResult, error) {
	id, err := m.createFn(ex)
	if err != nil {
		return nil, err
	}
	return &mongo.InsertOneResult{InsertedID: id}, nil
}
func (m *mockRepo) UpdateExercise(ex models.Exercise) (*mongo.UpdateResult, error) {
	ok, err := m.updateFn(ex)
	if err != nil {
		return nil, err
	}
	if ok {
		return &mongo.UpdateResult{MatchedCount: 1}, nil
	}
	return &mongo.UpdateResult{MatchedCount: 0}, nil
}
func (m *mockRepo) DeleteExercise(id primitive.ObjectID) (*mongo.DeleteResult, error) {
	ok, err := m.deleteFn(id)
	if err != nil {
		return nil, err
	}
	if ok {
		return &mongo.DeleteResult{DeletedCount: 1}, nil
	}
	return &mongo.DeleteResult{DeletedCount: 0}, nil
}

func TestGetExerciseByID_EmptyID(t *testing.T) {
	svc := NewExerciseService(&mockRepo{})
	_, err := svc.GetExerciseByID("")
	if err == nil {
		t.Fatalf("expected error when id is empty")
	}
}

func TestCreateExercise_Success(t *testing.T) {

	req := dto.ExerciseRequest{
		UserID:      primitive.NewObjectID().Hex(),
		Name:        "Push Up",
		Category:    "strength",
		MuscleGroup: "chest",
		Difficulty:  "easy",
	}

	createdModel := models.Exercise{
		ID:          primitive.NewObjectID(),
		UserID:      req.UserID,
		Name:        req.Name,
		Category:    req.Category,
		MuscleGroup: req.MuscleGroup,
		Difficulty:  req.Difficulty,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	repo := &mockRepo{
		createFn: func(ex models.Exercise) (primitive.ObjectID, error) {
			return createdModel.ID, nil
		},
		getByIDFn: func(id string) (models.Exercise, error) {
			return createdModel, nil
		},
	}

	svc := NewExerciseService(repo)
	res, err := svc.CreateExercise(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Name != req.Name || res.UserID != req.UserID {
		t.Fatalf("unexpected response: %+v", res)
	}
}

func TestUpdateExercise_Unauthorized(t *testing.T) {

	existing := models.Exercise{
		ID:     primitive.NewObjectID(),
		UserID: primitive.NewObjectID().Hex(),
	}
	req := dto.ExerciseRequest{UserID: primitive.NewObjectID().Hex(), Name: "X", Category: "c", MuscleGroup: "m", Difficulty: "d"}

	repo := &mockRepo{
		getByIDFn: func(id string) (models.Exercise, error) { return existing, nil },
	}
	svc := NewExerciseService(repo)
	_, err := svc.UpdateExercise(existing.ID.Hex(), req)
	if err == nil {
		t.Fatalf("expected unauthorized error")
	}
}

func TestDeleteExercise_Unauthorized(t *testing.T) {
	existing := models.Exercise{ID: primitive.NewObjectID(), UserID: primitive.NewObjectID().Hex()}
	repo := &mockRepo{
		getByIDFn: func(id string) (models.Exercise, error) { return existing, nil },
	}
	svc := NewExerciseService(repo)
	err := svc.DeleteExercise(primitive.NewObjectID().Hex(), existing.ID.Hex())
	if err == nil {
		t.Fatalf("expected unauthorized error on delete")
	}
}

func TestSearchExercises_Mapping(t *testing.T) {
	exercises := []models.Exercise{
		{ID: primitive.NewObjectID(), UserID: "u1", Name: "A", Category: "cat", MuscleGroup: "mg"},
		{ID: primitive.NewObjectID(), UserID: "u2", Name: "B", Category: "cat2", MuscleGroup: "mg2"},
	}
	repo := &mockRepo{
		getExercisesFn: func(name, category, muscleGroup string) ([]models.Exercise, error) { return exercises, nil },
	}
	svc := NewExerciseService(repo)
	res, err := svc.SearchExercises(dto.ExerciseSearch{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res) != len(exercises) {
		t.Fatalf("expected %d results, got %d", len(exercises), len(res))
	}
}
