package services

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"backend/dto"
	"backend/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mockRoutineRepo struct {
	created models.Routine
	store   map[string]models.Routine
}

func (m *mockRoutineRepo) GetRoutines(ownerID primitive.ObjectID, name string) ([]models.Routine, error) {
	out := []models.Routine{}
	for _, r := range m.store {
		if r.OwnerID == ownerID {
			if name == "" || (name != "" && r.Name == name) {
				out = append(out, r)
			}
		}
	}
	return out, nil
}

func (m *mockRoutineRepo) GetRoutineByID(id string) (models.Routine, error) {
	if r, ok := m.store[id]; ok {
		return r, nil
	}
	return models.Routine{}, errors.New("not found")
}

func (m *mockRoutineRepo) CreateRoutine(routine models.Routine) (*mongo.InsertOneResult, error) {
	m.created = routine
	if m.store == nil {
		m.store = map[string]models.Routine{}
	}
	m.store[routine.ID.Hex()] = routine

	return &mongo.InsertOneResult{InsertedID: routine.ID}, nil
}

func (m *mockRoutineRepo) UpdateRoutine(routine models.Routine) (*mongo.UpdateResult, error) {
	if m.store == nil {
		return nil, errors.New("no store")
	}
	m.store[routine.ID.Hex()] = routine
	return &mongo.UpdateResult{}, nil
}

func (m *mockRoutineRepo) DeleteRoutine(id primitive.ObjectID) (*mongo.DeleteResult, error) {
	if m.store == nil {
		return nil, errors.New("no store")
	}
	delete(m.store, id.Hex())
	return &mongo.DeleteResult{}, nil
}

type mockExerciseRepo struct {
	byID map[string]models.Exercise
}

func (m *mockExerciseRepo) GetExercises(name, category, muscleGroup string) ([]models.Exercise, error) {
	out := []models.Exercise{}
	for _, e := range m.byID {
		out = append(out, e)
	}
	return out, nil
}

func (m *mockExerciseRepo) GetExerciseByID(id string) (models.Exercise, error) {
	if e, ok := m.byID[id]; ok {
		return e, nil
	}
	return models.Exercise{}, nil
}

func (m *mockExerciseRepo) CreateExercise(exercise models.Exercise) (*mongo.InsertOneResult, error) {
	return &mongo.InsertOneResult{InsertedID: exercise.ID}, nil
}

func (m *mockExerciseRepo) UpdateExercise(exercise models.Exercise) (*mongo.UpdateResult, error) {
	return &mongo.UpdateResult{}, nil
}

func (m *mockExerciseRepo) DeleteExercise(id primitive.ObjectID) (*mongo.DeleteResult, error) {
	return &mongo.DeleteResult{}, nil
}

func TestCreateRoutine_Success(t *testing.T) {

	exerciseID := primitive.NewObjectID()
	ownerID := primitive.NewObjectID()

	exerRepo := &mockExerciseRepo{byID: map[string]models.Exercise{exerciseID.Hex(): {ID: exerciseID}}}
	repo := &mockRoutineRepo{store: map[string]models.Routine{}}
	svc := NewRoutineService(repo, exerRepo)

	req := dto.RoutineRequest{
		Name: "legs",
		Excercises: []dto.RoutineExcerciseList{{
			ExerciseID: exerciseID.Hex(),
			Order:      1,
			Sets:       3,
			Reps:       10,
		}},
	}

	got, err := svc.CreateRoutine(ownerID.Hex(), req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if got.Name != req.Name {
		t.Fatalf("expected name %s, got %s", req.Name, got.Name)
	}
	if len(got.Excercises) != 1 {
		t.Fatalf("expected 1 exercise, got %d", len(got.Excercises))
	}
}

func TestGetRoutines_Success(t *testing.T) {
	ownerID := primitive.NewObjectID()
	r := models.Routine{ID: primitive.NewObjectID(), OwnerID: ownerID, Name: "push", Entries: []models.RoutineExcerciseList{}}
	repo := &mockRoutineRepo{store: map[string]models.Routine{r.ID.Hex(): r}}
	svc := NewRoutineService(repo, &mockExerciseRepo{})

	out, err := svc.GetRoutines(ownerID.Hex(), "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 1 || out[0].Name != "push" {
		t.Fatalf("unexpected result: %#v", out)
	}
}

func TestGetRoutineByID_NotFound(t *testing.T) {
	repo := &mockRoutineRepo{store: map[string]models.Routine{}}
	svc := NewRoutineService(repo, &mockExerciseRepo{})

	_, err := svc.GetRoutineByID(primitive.NewObjectID().Hex())
	if err == nil {
		t.Fatalf("expected error for missing routine")
	}
}

func TestUpdateRoutine_Unauthorized(t *testing.T) {
	ownerID := primitive.NewObjectID()
	otherOwner := primitive.NewObjectID()
	r := models.Routine{ID: primitive.NewObjectID(), OwnerID: otherOwner, Name: "r1", Entries: []models.RoutineExcerciseList{{}}}
	repo := &mockRoutineRepo{store: map[string]models.Routine{r.ID.Hex(): r}}
	svc := NewRoutineService(repo, &mockExerciseRepo{})

	req := dto.RoutineRequest{Name: "new", Excercises: []dto.RoutineExcerciseList{{ExerciseID: primitive.NewObjectID().Hex(), Order: 1, Sets: 1, Reps: 1}}}
	_, err := svc.UpdateRoutine(ownerID.Hex(), r.ID.Hex(), req)
	if err == nil {
		t.Fatalf("expected unauthorized error")
	}
}

func TestDuplicateRoutine_Success(t *testing.T) {
	ownerID := primitive.NewObjectID()
	srcOwner := primitive.NewObjectID()
	exID := primitive.NewObjectID()
	src := models.Routine{
		ID:        primitive.NewObjectID(),
		OwnerID:   srcOwner,
		Name:      "orig",
		Entries:   []models.RoutineExcerciseList{{ExerciseID: exID, Order: 1, Sets: 2, Reps: 5}},
		IsPublic:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	exerRepo := &mockExerciseRepo{byID: map[string]models.Exercise{exID.Hex(): {ID: exID}}}
	repo := &mockRoutineRepo{store: map[string]models.Routine{src.ID.Hex(): src}}
	svc := NewRoutineService(repo, exerRepo)

	newName := "copy"
	idHex, err := svc.DuplicateRoutine(ownerID.Hex(), src.ID.Hex(), newName)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if idHex == "" {
		t.Fatalf("expected inserted id, got empty string")
	}

	stored, ok := repo.store[idHex]
	if !ok {
		t.Fatalf("expected stored copy under id %s", idHex)
	}
	if stored.Name != newName {
		t.Fatalf("expected name %s, got %s", newName, stored.Name)
	}
}

func TestValidateRoutineEntries_Errors(t *testing.T) {
	cases := []struct {
		name    string
		entries []dto.RoutineExcerciseList
		wantErr bool
	}{
		{"empty", []dto.RoutineExcerciseList{}, true},
		{"missing id", []dto.RoutineExcerciseList{{ExerciseID: "", Order: 1, Sets: 1, Reps: 1}}, true},
		{"invalid id", []dto.RoutineExcerciseList{{ExerciseID: "nothex", Order: 1, Sets: 1, Reps: 1}}, true},
		{"dup order", []dto.RoutineExcerciseList{{ExerciseID: primitive.NewObjectID().Hex(), Order: 1, Sets: 1, Reps: 1}, {ExerciseID: primitive.NewObjectID().Hex(), Order: 1, Sets: 1, Reps: 1}}, true},
		{"ok", []dto.RoutineExcerciseList{{ExerciseID: primitive.NewObjectID().Hex(), Order: 1, Sets: 1, Reps: 1}}, false},
	}

	for _, c := range cases {
		err := validateRoutineEntries(c.entries)
		if (err != nil) != c.wantErr {
			t.Fatalf("case %s: expected error=%v got %v", c.name, c.wantErr, err)
		}
	}
}

var _ = reflect.TypeOf((*mockRoutineRepo)(nil))
var _ = reflect.TypeOf((*mockExerciseRepo)(nil))
