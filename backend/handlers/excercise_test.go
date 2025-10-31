package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"backend/dto"
	"backend/models"
)

type mockExerciseService struct {
	getByIDFn func(id string) (models.Exercise, error)
	getListFn func(name, category, muscleGroup string) ([]models.Exercise, error)
	createFn  func(req dto.ExerciseRequest) (dto.ExerciseResponse, error)
	updateFn  func(id string, req dto.ExerciseRequest) (dto.ExerciseResponse, error)
	deleteFn  func(ownerID, exerciseID string) error
	searchFn  func(search dto.ExerciseSearch) ([]dto.ExerciseResponse, error)
}

func (m *mockExerciseService) GetExercises(name, category, muscleGroup string) ([]models.Exercise, error) {
	if m.getListFn != nil {
		return m.getListFn(name, category, muscleGroup)
	}
	return nil, nil
}
func (m *mockExerciseService) GetExerciseByID(id string) (models.Exercise, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(id)
	}
	return models.Exercise{}, nil
}
func (m *mockExerciseService) CreateExercise(req dto.ExerciseRequest) (dto.ExerciseResponse, error) {
	if m.createFn != nil {
		return m.createFn(req)
	}
	return dto.ExerciseResponse{}, nil
}
func (m *mockExerciseService) UpdateExercise(id string, req dto.ExerciseRequest) (dto.ExerciseResponse, error) {
	if m.updateFn != nil {
		return m.updateFn(id, req)
	}
	return dto.ExerciseResponse{}, nil
}
func (m *mockExerciseService) DeleteExercise(ownerID, exerciseID string) error {
	if m.deleteFn != nil {
		return m.deleteFn(ownerID, exerciseID)
	}
	return nil
}
func (m *mockExerciseService) SearchExercises(search dto.ExerciseSearch) ([]dto.ExerciseResponse, error) {
	if m.searchFn != nil {
		return m.searchFn(search)
	}
	return nil, nil
}

func makeReqWithCtx(t *testing.T, method, path string, body interface{}, ctxVals map[string]interface{}) (*gin.Context, *httptest.ResponseRecorder) {
	t.Helper()
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			t.Fatalf("encode body: %v", err)
		}
	}
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(method, path, &buf)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	c.Request = req
	for k, v := range ctxVals {
		c.Set(k, v)
	}
	return c, w
}

func TestGetExercise_ByID_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	now := time.Now()
	ex := models.Exercise{
		ID:          primitive.NewObjectID(),
		Name:        "Push Up",
		Category:    "strength",
		MuscleGroup: "chest",
		Difficulty:  "easy",
		UserID:      primitive.NewObjectID().Hex(),
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	svc := &mockExerciseService{
		getByIDFn: func(id string) (models.Exercise, error) {
			return ex, nil
		},
	}

	h := NewExerciseHandler(svc)

	c, w := makeReqWithCtx(t, "GET", "/exercises/123", nil, nil)

	c.Params = gin.Params{{Key: "id", Value: "123"}}

	h.GetExercise(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d got %d body=%s", http.StatusOK, w.Code, w.Body.String())
	}
}

func TestGetExercise_Search_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	now := time.Now()
	exList := []models.Exercise{{
		ID:          primitive.NewObjectID(),
		Name:        "Squat",
		Category:    "strength",
		MuscleGroup: "legs",
		Difficulty:  "medium",
		UserID:      primitive.NewObjectID().Hex(),
		CreatedAt:   now,
		UpdatedAt:   now,
	}}

	svc := &mockExerciseService{
		getListFn: func(name, category, muscleGroup string) ([]models.Exercise, error) {
			return exList, nil
		},
	}

	h := NewExerciseHandler(svc)

	c, w := makeReqWithCtx(t, "GET", "/exercises", nil, nil)

	q := c.Request.URL.Query()
	q.Add("name", "Squat")
	c.Request.URL.RawQuery = q.Encode()

	h.GetExercise(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d got %d body=%s", http.StatusOK, w.Code, w.Body.String())
	}
}

func TestCreateExercise_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svc := &mockExerciseService{}
	h := NewExerciseHandler(svc)

	req := dto.ExerciseRequest{Name: "X", Category: "c", MuscleGroup: "m", Difficulty: "d"}

	c, w := makeReqWithCtx(t, "POST", "/exercises", req, nil)

	h.CreateExercise(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d got %d body=%s", http.StatusUnauthorized, w.Code, w.Body.String())
	}
}

func TestCreateExercise_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	resp := dto.ExerciseResponse{ID: primitive.NewObjectID().Hex(), Name: "Bench"}
	svc := &mockExerciseService{
		createFn: func(req dto.ExerciseRequest) (dto.ExerciseResponse, error) {
			return resp, nil
		},
	}

	h := NewExerciseHandler(svc)

	ctxVals := map[string]interface{}{"user_id": primitive.NewObjectID().Hex(), "user_role": "admin"}
	req := dto.ExerciseRequest{Name: "Bench", Category: "c", MuscleGroup: "m", Difficulty: "d"}
	c, w := makeReqWithCtx(t, "POST", "/exercises", req, ctxVals)

	h.CreateExercise(c)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status %d got %d body=%s", http.StatusCreated, w.Code, w.Body.String())
	}
}

func TestUpdateExercise_BadRequest_NoID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svc := &mockExerciseService{}
	h := NewExerciseHandler(svc)
	ctxVals := map[string]interface{}{"user_id": primitive.NewObjectID().Hex(), "user_role": "admin"}
	req := dto.ExerciseRequest{Name: "Bench", Category: "c", MuscleGroup: "m", Difficulty: "d"}
	c, w := makeReqWithCtx(t, "PUT", "/exercises/", req, ctxVals)

	h.UpdateExercise(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d got %d body=%s", http.StatusBadRequest, w.Code, w.Body.String())
	}
}

func TestDeleteExercise_MissingID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svc := &mockExerciseService{}
	h := NewExerciseHandler(svc)
	ctxVals := map[string]interface{}{"user_id": primitive.NewObjectID().Hex(), "user_role": "admin"}
	c, w := makeReqWithCtx(t, "DELETE", "/exercises/", nil, ctxVals)

	h.DeleteExercise(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d got %d body=%s", http.StatusBadRequest, w.Code, w.Body.String())
	}
}
