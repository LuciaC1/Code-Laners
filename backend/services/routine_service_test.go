package services

import (
	"errors"
	"testing"
	"time"

	"backend/dto"
	"backend/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Mock repositories
type MockRoutineRepository struct {
	mock.Mock
}

type MockExerciseRepository struct {
	mock.Mock
}

// Implement GetExercises to satisfy ExerciseRepositoryInterface
func (m *MockExerciseRepository) GetExercises(ownerID string, name string, filter string) ([]models.Exercise, error) {
	args := m.Called(ownerID, name, filter)
	return args.Get(0).([]models.Exercise), args.Error(1)
}

// Implement RoutineRepositoryInterface
func (m *MockRoutineRepository) CreateRoutine(routine models.Routine) (*mongo.InsertOneResult, error) {
	args := m.Called(routine)
	return args.Get(0).(*mongo.InsertOneResult), args.Error(1)
}

func (m *MockRoutineRepository) GetRoutines(ownerID primitive.ObjectID, name string) ([]models.Routine, error) {
	args := m.Called(ownerID, name)
	return args.Get(0).([]models.Routine), args.Error(1)
}

func (m *MockRoutineRepository) GetRoutineByID(id string) (models.Routine, error) {
	args := m.Called(id)
	return args.Get(0).(models.Routine), args.Error(1)
}

func (m *MockRoutineRepository) UpdateRoutine(routine models.Routine) (*mongo.UpdateResult, error) {
	args := m.Called(routine)
	return args.Get(0).(*mongo.UpdateResult), args.Error(1)
}

func (m *MockRoutineRepository) DeleteRoutine(id primitive.ObjectID) (*mongo.DeleteResult, error) {
	args := m.Called(id)
	return args.Get(0).(*mongo.DeleteResult), args.Error(1)
}

// Implement ExerciseRepositoryInterface
func (m *MockExerciseRepository) GetExerciseByID(id string) (models.Exercise, error) {
	args := m.Called(id)
	return args.Get(0).(models.Exercise), args.Error(1)
}

func (m *MockExerciseRepository) CreateExercise(exercise models.Exercise) (*mongo.InsertOneResult, error) {
	args := m.Called(exercise)
	return args.Get(0).(*mongo.InsertOneResult), args.Error(1)
}

// Implement DeleteExercise to satisfy ExerciseRepositoryInterface
func (m *MockExerciseRepository) DeleteExercise(id primitive.ObjectID) (*mongo.DeleteResult, error) {
	args := m.Called(id)
	return args.Get(0).(*mongo.DeleteResult), args.Error(1)
}

// Implement UpdateExercise to satisfy ExerciseRepositoryInterface
func (m *MockExerciseRepository) UpdateExercise(exercise models.Exercise) (*mongo.UpdateResult, error) {
	args := m.Called(exercise)
	return args.Get(0).(*mongo.UpdateResult), args.Error(1)
}

func TestCreateRoutine(t *testing.T) {
	mockRoutineRepo := new(MockRoutineRepository)
	mockExerciseRepo := new(MockExerciseRepository)
	service := NewRoutineService(mockRoutineRepo, mockExerciseRepo)

	t.Run("Success", func(t *testing.T) {
		ownerID := primitive.NewObjectID().Hex()
		exerciseID := primitive.NewObjectID().Hex()

		input := dto.RoutineRequest{
			Name: "Test Routine",
			Excercises: []dto.RoutineExcerciseList{
				{
					ExerciseID: exerciseID,
					Order:      1,
					Sets:       3,
					Reps:       10,
					Weight:     20.5,
				},
			},
		}

		exID, err := primitive.ObjectIDFromHex(exerciseID)
		assert.NoError(t, err)
		mockExerciseRepo.On("GetExerciseByID", exerciseID).Return(models.Exercise{
			ID: exID,
		}, nil)

		mockRoutineRepo.On("CreateRoutine", mock.AnythingOfType("models.Routine")).Return(
			&mongo.InsertOneResult{InsertedID: primitive.NewObjectID()},
			nil,
		)

		result, err := service.CreateRoutine(ownerID, input)

		assert.NoError(t, err)
		assert.NotEmpty(t, result)
		assert.Equal(t, input.Name, result.Name)
	})

	t.Run("Invalid OwnerID", func(t *testing.T) {
		_, err := service.CreateRoutine("invalid-id", dto.RoutineRequest{})
		assert.Error(t, err)
	})

	t.Run("Empty Exercises", func(t *testing.T) {
		ownerID := primitive.NewObjectID().Hex()
		_, err := service.CreateRoutine(ownerID, dto.RoutineRequest{
			Name:       "Test",
			Excercises: []dto.RoutineExcerciseList{},
		})
		assert.Error(t, err)
	})
}

func TestGetRoutines(t *testing.T) {
	mockRoutineRepo := new(MockRoutineRepository)
	mockExerciseRepo := new(MockExerciseRepository)
	service := NewRoutineService(mockRoutineRepo, mockExerciseRepo)

	t.Run("Success", func(t *testing.T) {
		ownerID := primitive.NewObjectID().Hex()
		routines := []models.Routine{
			{
				ID:        primitive.NewObjectID(),
				Name:      "Routine 1",
				CreatedAt: time.Now(),
			},
		}

		mockRoutineRepo.On("GetRoutines", mock.AnythingOfType("primitive.ObjectID"), "").Return(routines, nil)

		result, err := service.GetRoutines(ownerID, "")

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, routines[0].Name, result[0].Name)
	})

	t.Run("Invalid OwnerID", func(t *testing.T) {
		_, err := service.GetRoutines("invalid-id", "")
		assert.Error(t, err)
	})
}

func TestUpdateRoutine(t *testing.T) {
	mockRoutineRepo := new(MockRoutineRepository)
	mockExerciseRepo := new(MockExerciseRepository)
	service := NewRoutineService(mockRoutineRepo, mockExerciseRepo)

	t.Run("Success", func(t *testing.T) {
		ownerID := primitive.NewObjectID()
		routineID := primitive.NewObjectID().Hex()
		exerciseID := primitive.NewObjectID().Hex()

		routineObjID, err := primitive.ObjectIDFromHex(routineID)
		assert.NoError(t, err)
		existingRoutine := models.Routine{
			ID:      routineObjID,
			OwnerID: ownerID,
		}

		input := dto.RoutineRequest{
			Name: "Updated Routine",
			Excercises: []dto.RoutineExcerciseList{
				{
					ExerciseID: exerciseID,
					Order:      1,
					Sets:       3,
					Reps:       10,
				},
			},
		}

		mockRoutineRepo.On("GetRoutineByID", routineID).Return(existingRoutine, nil)
		exID, err := primitive.ObjectIDFromHex(exerciseID)
		assert.NoError(t, err)
		mockExerciseRepo.On("GetExerciseByID", exerciseID).Return(models.Exercise{
			ID: exID,
		}, nil)
		mockRoutineRepo.On("UpdateRoutine", mock.AnythingOfType("models.Routine")).Return(&mongo.UpdateResult{}, nil)

		result, err := service.UpdateRoutine(ownerID.Hex(), routineID, input)

		assert.NoError(t, err)
		assert.Equal(t, input.Name, result.Name)
	})

	t.Run("Not Owner", func(t *testing.T) {
		ownerID := primitive.NewObjectID()
		differentOwnerID := primitive.NewObjectID()
		routineID := primitive.NewObjectID().Hex()

		routineObjID, err := primitive.ObjectIDFromHex(routineID)
		assert.NoError(t, err)
		existingRoutine := models.Routine{
			ID:      routineObjID,
			OwnerID: ownerID,
		}

		mockRoutineRepo.On("GetRoutineByID", routineID).Return(existingRoutine, nil)

		_, err = service.UpdateRoutine(differentOwnerID.Hex(), routineID, dto.RoutineRequest{})

		assert.Error(t, err)
	})
}

func TestDeleteRoutine(t *testing.T) {
	mockRoutineRepo := new(MockRoutineRepository)
	mockExerciseRepo := new(MockExerciseRepository)
	service := NewRoutineService(mockRoutineRepo, mockExerciseRepo)

	t.Run("Success", func(t *testing.T) {
		ownerID := primitive.NewObjectID()
		routineID := primitive.NewObjectID().Hex()

		routineObjID, err := primitive.ObjectIDFromHex(routineID)
		assert.NoError(t, err)
		existingRoutine := models.Routine{
			ID:      routineObjID,
			OwnerID: ownerID,
		}

		mockRoutineRepo.On("GetRoutineByID", routineID).Return(existingRoutine, nil)
		mockRoutineRepo.On("DeleteRoutine", existingRoutine.ID).Return(&mongo.DeleteResult{DeletedCount: 1}, nil)

		err = service.DeleteRoutine(ownerID.Hex(), routineID)
		assert.NoError(t, err)
	})

	t.Run("Not Found", func(t *testing.T) {
		routineID := primitive.NewObjectID().Hex()
		mockRoutineRepo.On("GetRoutineByID", routineID).Return(models.Routine{}, errors.New("not found"))

		err := service.DeleteRoutine(primitive.NewObjectID().Hex(), routineID)
		assert.Error(t, err)
	})
}

func TestDuplicateRoutine(t *testing.T) {
	mockRoutineRepo := new(MockRoutineRepository)
	mockExerciseRepo := new(MockExerciseRepository)
	service := NewRoutineService(mockRoutineRepo, mockExerciseRepo)

	t.Run("Success", func(t *testing.T) {
		ownerID := primitive.NewObjectID().Hex()
		sourceRoutineID := primitive.NewObjectID().Hex()
		exerciseID := primitive.NewObjectID()

		sourceRoutineObjID, err := primitive.ObjectIDFromHex(sourceRoutineID)
		assert.NoError(t, err)
		sourceRoutine := models.Routine{
			ID:          sourceRoutineObjID,
			Name:        "Source Routine",
			Description: "Test Description",
			Entries: []models.RoutineExcerciseList{
				{
					ExerciseID: exerciseID,
					Order:      1,
					Sets:       3,
					Reps:       10,
				},
			},
		}

		mockRoutineRepo.On("GetRoutineByID", sourceRoutineID).Return(sourceRoutine, nil)
		mockExerciseRepo.On("GetExerciseByID", exerciseID.Hex()).Return(models.Exercise{ID: exerciseID}, nil)
		mockRoutineRepo.On("CreateRoutine", mock.AnythingOfType("models.Routine")).Return(
			&mongo.InsertOneResult{InsertedID: primitive.NewObjectID()},
			nil,
		)

		newID, err := service.DuplicateRoutine(ownerID, sourceRoutineID, "New Routine")

		assert.NoError(t, err)
		assert.NotEmpty(t, newID)
	})

	t.Run("Source Not Found", func(t *testing.T) {
		sourceRoutineID := primitive.NewObjectID().Hex()
		mockRoutineRepo.On("GetRoutineByID", sourceRoutineID).Return(models.Routine{}, errors.New("not found"))

		_, err := service.DuplicateRoutine(primitive.NewObjectID().Hex(), sourceRoutineID, "New Name")
		assert.Error(t, err)
	})
}

func TestGetRoutineByID_SuccessAndError(t *testing.T) {
	mockRoutineRepo := new(MockRoutineRepository)
	mockExerciseRepo := new(MockExerciseRepository)
	service := NewRoutineService(mockRoutineRepo, mockExerciseRepo)

	t.Run("Success", func(t *testing.T) {
		routineID := primitive.NewObjectID().Hex()
		objID, _ := primitive.ObjectIDFromHex(routineID)
		model := models.Routine{
			ID:        objID,
			Name:      "My Routine",
			CreatedAt: time.Now(),
		}
		mockRoutineRepo.On("GetRoutineByID", routineID).Return(model, nil)

		res, err := service.GetRoutineByID(routineID)
		assert.NoError(t, err)
		assert.Equal(t, model.Name, res.Name)
	})

	t.Run("NotFound", func(t *testing.T) {
		routineID := primitive.NewObjectID().Hex()
		mockRoutineRepo.On("GetRoutineByID", routineID).Return(models.Routine{}, errors.New("not found"))

		_, err := service.GetRoutineByID(routineID)
		assert.Error(t, err)
	})
}

func TestCreateRoutine_ValidationAndMissingExercise(t *testing.T) {
	mockRoutineRepo := new(MockRoutineRepository)
	mockExerciseRepo := new(MockExerciseRepository)
	service := NewRoutineService(mockRoutineRepo, mockExerciseRepo)

	t.Run("Missing Exercise In Repo", func(t *testing.T) {
		ownerID := primitive.NewObjectID().Hex()
		exerciseID := primitive.NewObjectID().Hex()

		input := dto.RoutineRequest{
			Name: "Routine Missing Exercise",
			Excercises: []dto.RoutineExcerciseList{
				{ExerciseID: exerciseID, Order: 1, Sets: 3, Reps: 8},
			},
		}

		mockExerciseRepo.On("GetExerciseByID", exerciseID).Return(models.Exercise{}, nil)

		_, err := service.CreateRoutine(ownerID, input)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "exercises not found")
	})

	t.Run("Invalid ExerciseID in Entry", func(t *testing.T) {
		ownerID := primitive.NewObjectID().Hex()

		input := dto.RoutineRequest{
			Name: "Routine Invalid Entry",
			Excercises: []dto.RoutineExcerciseList{
				{ExerciseID: "invalid-hex", Order: 1, Sets: 3, Reps: 8},
			},
		}
		_, err := service.CreateRoutine(ownerID, input)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "exercise_id inválido")
	})

	t.Run("Duplicate Order in Entries", func(t *testing.T) {
		ownerID := primitive.NewObjectID().Hex()
		ex1 := primitive.NewObjectID().Hex()
		ex2 := primitive.NewObjectID().Hex()

		input := dto.RoutineRequest{
			Name: "Routine Duplicate Order",
			Excercises: []dto.RoutineExcerciseList{
				{ExerciseID: ex1, Order: 1, Sets: 3, Reps: 8},
				{ExerciseID: ex2, Order: 1, Sets: 4, Reps: 10},
			},
		}

		_, err := service.CreateRoutine(ownerID, input)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "order duplicado")
	})

	t.Run("Non positive sets or reps", func(t *testing.T) {
		ownerID := primitive.NewObjectID().Hex()
		ex := primitive.NewObjectID().Hex()

		inputSetsZero := dto.RoutineRequest{
			Name: "Sets Zero",
			Excercises: []dto.RoutineExcerciseList{
				{ExerciseID: ex, Order: 1, Sets: 0, Reps: 8},
			},
		}
		_, err := service.CreateRoutine(ownerID, inputSetsZero)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "sets debe ser > 0")

		inputRepsZero := dto.RoutineRequest{
			Name: "Reps Zero",
			Excercises: []dto.RoutineExcerciseList{
				{ExerciseID: ex, Order: 1, Sets: 3, Reps: 0},
			},
		}
		_, err = service.CreateRoutine(ownerID, inputRepsZero)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "reps debe ser > 0")
	})
}

func TestGetRoutines_ErrorFromRepo(t *testing.T) {
	mockRoutineRepo := new(MockRoutineRepository)
	mockExerciseRepo := new(MockExerciseRepository)
	service := NewRoutineService(mockRoutineRepo, mockExerciseRepo)

	ownerID := primitive.NewObjectID().Hex()
	mockRoutineRepo.On("GetRoutines", mock.AnythingOfType("primitive.ObjectID"), "").Return(nil, errors.New("db error"))

	_, err := service.GetRoutines(ownerID, "")
	assert.Error(t, err)
}

func TestUpdateRoutine_Errors(t *testing.T) {
	mockRoutineRepo := new(MockRoutineRepository)
	mockExerciseRepo := new(MockExerciseRepository)
	service := NewRoutineService(mockRoutineRepo, mockExerciseRepo)

	t.Run("Invalid ownerID format", func(t *testing.T) {
		routineID := primitive.NewObjectID().Hex()
		routineObjID, _ := primitive.ObjectIDFromHex(routineID)
		existing := models.Routine{ID: routineObjID, OwnerID: primitive.NewObjectID()}

		mockRoutineRepo.On("GetRoutineByID", routineID).Return(existing, nil)

		_, err := service.UpdateRoutine("invalid-owner", routineID, dto.RoutineRequest{})
		assert.Error(t, err)
	})

	t.Run("Owner not authorized", func(t *testing.T) {
		owner := primitive.NewObjectID()
		otherOwner := primitive.NewObjectID()
		routineID := primitive.NewObjectID().Hex()
		routineObjID, _ := primitive.ObjectIDFromHex(routineID)
		existing := models.Routine{ID: routineObjID, OwnerID: owner}

		mockRoutineRepo.On("GetRoutineByID", routineID).Return(existing, nil)

		_, err := service.UpdateRoutine(otherOwner.Hex(), routineID, dto.RoutineRequest{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no autorizado")
	})
}

func TestDeleteRoutine_Errors(t *testing.T) {
	mockRoutineRepo := new(MockRoutineRepository)
	mockExerciseRepo := new(MockExerciseRepository)
	service := NewRoutineService(mockRoutineRepo, mockExerciseRepo)

	t.Run("Owner invalid format", func(t *testing.T) {
		routineID := primitive.NewObjectID().Hex()
		routineObjID, _ := primitive.ObjectIDFromHex(routineID)
		existing := models.Routine{ID: routineObjID, OwnerID: primitive.NewObjectID()}

		mockRoutineRepo.On("GetRoutineByID", routineID).Return(existing, nil)

		err := service.DeleteRoutine("invalid-owner", routineID)
		assert.Error(t, err)
	})

	t.Run("Not owner", func(t *testing.T) {
		owner := primitive.NewObjectID()
		routineID := primitive.NewObjectID().Hex()
		routineObjID, _ := primitive.ObjectIDFromHex(routineID)
		existing := models.Routine{ID: routineObjID, OwnerID: owner}

		mockRoutineRepo.On("GetRoutineByID", routineID).Return(existing, nil)

		// use a different owner
		err := service.DeleteRoutine(primitive.NewObjectID().Hex(), routineID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no autorizado")
	})
}

func TestDuplicateRoutine_ErrorsAndMissingExercises(t *testing.T) {
	mockRoutineRepo := new(MockRoutineRepository)
	mockExerciseRepo := new(MockExerciseRepository)
	service := NewRoutineService(mockRoutineRepo, mockExerciseRepo)

	t.Run("Invalid ownerID format", func(t *testing.T) {
		sourceID := primitive.NewObjectID().Hex()
		objID, _ := primitive.ObjectIDFromHex(sourceID)
		src := models.Routine{ID: objID, Entries: []models.RoutineExcerciseList{}}

		mockRoutineRepo.On("GetRoutineByID", sourceID).Return(src, nil)

		_, err := service.DuplicateRoutine("invalid-owner", sourceID, "copy")
		assert.Error(t, err)
	})

	t.Run("Missing exercise in source", func(t *testing.T) {
		sourceID := primitive.NewObjectID().Hex()
		exID := primitive.NewObjectID()
		objID, _ := primitive.ObjectIDFromHex(sourceID)
		src := models.Routine{
			ID: objID,
			Entries: []models.RoutineExcerciseList{
				{ExerciseID: exID, Order: 1, Sets: 3, Reps: 8},
			},
		}

		mockRoutineRepo.On("GetRoutineByID", sourceID).Return(src, nil)
		// repo returns empty exercise -> reported as missing
		mockExerciseRepo.On("GetExerciseByID", exID.Hex()).Return(models.Exercise{}, nil)

		_, err := service.DuplicateRoutine(primitive.NewObjectID().Hex(), sourceID, "copy")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "exercises not found")
	})

	t.Run("Error when creating copy", func(t *testing.T) {
		sourceID := primitive.NewObjectID().Hex()
		exID := primitive.NewObjectID()
		objID, _ := primitive.ObjectIDFromHex(sourceID)
		src := models.Routine{
			ID: objID,
			Entries: []models.RoutineExcerciseList{
				{ExerciseID: exID, Order: 1, Sets: 3, Reps: 8},
			},
		}

		mockRoutineRepo.On("GetRoutineByID", sourceID).Return(src, nil)
		mockExerciseRepo.On("GetExerciseByID", exID.Hex()).Return(models.Exercise{ID: exID}, nil)
		// simulate create error
		mockRoutineRepo.On("CreateRoutine", mock.AnythingOfType("models.Routine")).Return((*mongo.InsertOneResult)(nil), errors.New("create error"))

		_, err := service.DuplicateRoutine(primitive.NewObjectID().Hex(), sourceID, "copy")
		assert.Error(t, err)
	})
}
