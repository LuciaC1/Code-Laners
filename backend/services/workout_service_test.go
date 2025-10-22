package services

import (
	"errors"
	"testing"
	"time"

	"backend/dto"
	"backend/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// mockRepo implements repositories.WorkoutRepositoryInterface for tests
type mockWorkoutRepo struct {
	getWorkoutsFn    func(userID primitive.ObjectID) ([]models.Workout, error)
	getWorkoutByIDFn func(id string) (models.Workout, error)
	createWorkoutFn  func(workout models.Workout) (*mongo.InsertOneResult, error)
	updateWorkoutFn  func(workout models.Workout) (*mongo.UpdateResult, error)
	deleteWorkoutFn  func(id primitive.ObjectID) (*mongo.DeleteResult, error)
}

func (m *mockWorkoutRepo) GetWorkouts(userID primitive.ObjectID) ([]models.Workout, error) {
	return m.getWorkoutsFn(userID)
}
func (m *mockWorkoutRepo) GetWorkoutByID(id string) (models.Workout, error) {
	return m.getWorkoutByIDFn(id)
}
func (m *mockWorkoutRepo) CreateWorkout(workout models.Workout) (*mongo.InsertOneResult, error) {
	return m.createWorkoutFn(workout)
}
func (m *mockWorkoutRepo) UpdateWorkout(workout models.Workout) (*mongo.UpdateResult, error) {
	return m.updateWorkoutFn(workout)
}
func (m *mockWorkoutRepo) DeleteWorkout(id primitive.ObjectID) (*mongo.DeleteResult, error) {
	return m.deleteWorkoutFn(id)
}

func TestGetWorkouts_Success(t *testing.T) {
	uid := primitive.NewObjectID()
	now := time.Now()
	sample := models.Workout{ID: primitive.NewObjectID(), UserID: uid, UpdatedAt: now}

	repo := &mockWorkoutRepo{
		getWorkoutsFn: func(userID primitive.ObjectID) ([]models.Workout, error) {
			if userID != uid {
				t.Fatalf("unexpected userID")
			}
			return []models.Workout{sample}, nil
		},
	}

	svc := NewWorkoutService(repo)
	dtos, err := svc.GetWorkouts(uid.Hex())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(dtos) != 1 {
		t.Fatalf("expected 1 workout, got %d", len(dtos))
	}
	if dtos[0].UserID != uid.Hex() {
		t.Fatalf("unexpected userID in dto")
	}
}

func TestGetWorkouts_InvalidUserID(t *testing.T) {
	svc := NewWorkoutService(&mockWorkoutRepo{})
	_, err := svc.GetWorkouts("")
	if err == nil {
		t.Fatalf("expected error for empty userID")
	}
}

func TestGetWorkoutByID_Success(t *testing.T) {
	id := primitive.NewObjectID().Hex()
	sample := models.Workout{ID: primitive.NewObjectID(), UpdatedAt: time.Now()}

	repo := &mockWorkoutRepo{
		getWorkoutByIDFn: func(i string) (models.Workout, error) {
			if i != id {
				return models.Workout{}, errors.New("not found")
			}
			return sample, nil
		},
	}

	svc := NewWorkoutService(repo)
	_, err := svc.GetWorkoutByID(id)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestGetWorkoutByID_Error(t *testing.T) {
	repo := &mockWorkoutRepo{
		getWorkoutByIDFn: func(i string) (models.Workout, error) {
			return models.Workout{}, errors.New("db error")
		},
	}
	svc := NewWorkoutService(repo)
	_, err := svc.GetWorkoutByID("any")
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestCreateWorkout_Success(t *testing.T) {
	uid := primitive.NewObjectID()
	input := dto.WorkoutDTO{UserID: uid.Hex(), DurationMinutes: 30}

	repo := &mockWorkoutRepo{
		createWorkoutFn: func(w models.Workout) (*mongo.InsertOneResult, error) {
			return &mongo.InsertOneResult{InsertedID: w.ID}, nil
		},
	}

	svc := NewWorkoutService(repo)
	id, err := svc.CreateWorkout(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id == "" {
		t.Fatalf("expected inserted id")
	}
}

func TestCreateWorkout_MissingUser(t *testing.T) {
	svc := NewWorkoutService(&mockWorkoutRepo{})
	_, err := svc.CreateWorkout(dto.WorkoutDTO{})
	if err == nil {
		t.Fatalf("expected error when user_id missing")
	}
}

func TestUpdateWorkout_Success(t *testing.T) {
	uid := primitive.NewObjectID()
	wid := primitive.NewObjectID()
	input := dto.WorkoutDTO{ID: wid, UserID: uid.Hex(), DurationMinutes: 45}

	repo := &mockWorkoutRepo{
		updateWorkoutFn: func(w models.Workout) (*mongo.UpdateResult, error) {
			if w.ID != wid {
				t.Fatalf("unexpected id")
			}
			return &mongo.UpdateResult{MatchedCount: 1}, nil
		},
	}

	svc := NewWorkoutService(repo)
	err := svc.UpdateWorkout(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUpdateWorkout_MissingID(t *testing.T) {
	svc := NewWorkoutService(&mockWorkoutRepo{})
	err := svc.UpdateWorkout(dto.WorkoutDTO{})
	if err == nil {
		t.Fatalf("expected error when id missing")
	}
}

func TestDeleteWorkout_Success(t *testing.T) {
	id := primitive.NewObjectID()
	repo := &mockWorkoutRepo{
		deleteWorkoutFn: func(i primitive.ObjectID) (*mongo.DeleteResult, error) {
			if i != id {
				t.Fatalf("unexpected id")
			}
			return &mongo.DeleteResult{DeletedCount: 1}, nil
		},
	}

	svc := NewWorkoutService(repo)
	err := svc.DeleteWorkout(id.Hex())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeleteWorkout_InvalidID(t *testing.T) {
	svc := NewWorkoutService(&mockWorkoutRepo{})
	err := svc.DeleteWorkout("")
	if err == nil {
		t.Fatalf("expected error for empty id")
	}
}
