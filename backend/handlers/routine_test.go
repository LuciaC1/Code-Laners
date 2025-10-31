package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"backend/dto"

	"github.com/gin-gonic/gin"
)

type mockRoutineService struct {
	GetRoutineByIDFunc func(id string) (dto.RoutineResponse, error)
	GetRoutinesFunc    func(ownerID, name string) ([]dto.RoutineResponse, error)
	CreateRoutineFunc  func(ownerID string, input dto.RoutineRequest) (dto.RoutineResponse, error)
	UpdateRoutineFunc  func(ownerID, routineID string, input dto.RoutineRequest) (dto.RoutineResponse, error)
	DeleteRoutineFunc  func(ownerID, routineID string) error
	DuplicateFunc      func(ownerID, sourceRoutineID, newName string) (string, error)
}

func (m *mockRoutineService) CreateRoutine(ownerID string, input dto.RoutineRequest) (dto.RoutineResponse, error) {
	if m.CreateRoutineFunc != nil {
		return m.CreateRoutineFunc(ownerID, input)
	}
	return dto.RoutineResponse{}, nil
}
func (m *mockRoutineService) GetRoutines(ownerID string, name string) ([]dto.RoutineResponse, error) {
	if m.GetRoutinesFunc != nil {
		return m.GetRoutinesFunc(ownerID, name)
	}
	return nil, nil
}
func (m *mockRoutineService) GetRoutineByID(id string) (dto.RoutineResponse, error) {
	if m.GetRoutineByIDFunc != nil {
		return m.GetRoutineByIDFunc(id)
	}
	return dto.RoutineResponse{}, nil
}
func (m *mockRoutineService) UpdateRoutine(ownerID string, routineID string, input dto.RoutineRequest) (dto.RoutineResponse, error) {
	if m.UpdateRoutineFunc != nil {
		return m.UpdateRoutineFunc(ownerID, routineID, input)
	}
	return dto.RoutineResponse{}, nil
}
func (m *mockRoutineService) DeleteRoutine(ownerID string, routineID string) error {
	if m.DeleteRoutineFunc != nil {
		return m.DeleteRoutineFunc(ownerID, routineID)
	}
	return nil
}
func (m *mockRoutineService) DuplicateRoutine(ownerID string, sourceRoutineID string, newName string) (string, error) {
	if m.DuplicateFunc != nil {
		return m.DuplicateFunc(ownerID, sourceRoutineID, newName)
	}
	return "", nil
}

func setupGinTest() {
	gin.SetMode(gin.TestMode)
}

func TestGetRoutineByID_Success(t *testing.T) {
	setupGinTest()

	userID := "user123"
	routineID := "routine123"

	mock := &mockRoutineService{
		GetRoutineByIDFunc: func(id string) (dto.RoutineResponse, error) {
			return dto.RoutineResponse{ID: id, UserID: userID, Name: "r1", Excercises: []dto.RoutineExcerciseList{}}, nil
		},
	}

	handler := NewRoutineHandler(mock)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: routineID}}
	c.Set("user_id", userID)

	handler.GetRoutineByID(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d, body: %s", w.Code, w.Body.String())
	}

	var got dto.RoutineResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if got.ID != routineID || got.UserID != userID {
		t.Fatalf("unexpected routine in response: %+v", got)
	}
}

func TestGetRoutineByID_Unauthorized(t *testing.T) {
	setupGinTest()
	mock := &mockRoutineService{}
	handler := NewRoutineHandler(mock)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "any"}}

	handler.GetRoutineByID(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", w.Code)
	}
}

func TestGetRoutineByID_NotFoundAndForbidden(t *testing.T) {
	setupGinTest()

	mock1 := &mockRoutineService{
		GetRoutineByIDFunc: func(id string) (dto.RoutineResponse, error) {
			return dto.RoutineResponse{}, errors.New("not found")
		},
	}
	handler1 := NewRoutineHandler(mock1)
	w1 := httptest.NewRecorder()
	c1, _ := gin.CreateTestContext(w1)
	c1.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c1.Params = gin.Params{{Key: "id", Value: "r"}}
	c1.Set("user_id", "u")
	handler1.GetRoutineByID(c1)
	if w1.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w1.Code)
	}

	mock2 := &mockRoutineService{
		GetRoutineByIDFunc: func(id string) (dto.RoutineResponse, error) {
			return dto.RoutineResponse{ID: id, UserID: "other", IsPublic: false}, nil
		},
	}
	handler2 := NewRoutineHandler(mock2)
	w2 := httptest.NewRecorder()
	c2, _ := gin.CreateTestContext(w2)
	c2.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c2.Params = gin.Params{{Key: "id", Value: "r"}}
	c2.Set("user_id", "u")
	handler2.GetRoutineByID(c2)
	if w2.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", w2.Code)
	}
}

func TestGetRoutines_SuccessAndUnauthorized(t *testing.T) {
	setupGinTest()

	mock := &mockRoutineService{
		GetRoutinesFunc: func(ownerID, name string) ([]dto.RoutineResponse, error) {
			return []dto.RoutineResponse{{ID: "r1", UserID: ownerID, Name: "n1"}}, nil
		},
	}
	handler := NewRoutineHandler(mock)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(http.MethodGet, "/?name=x", nil)
	c.Request = req
	c.Set("user_id", "u1")
	handler.GetRoutines(c)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body: %s", w.Code, w.Body.String())
	}

	w2 := httptest.NewRecorder()
	c2, _ := gin.CreateTestContext(w2)
	c2.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	handler.GetRoutines(c2)
	if w2.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w2.Code)
	}
}

func TestCreateRoutine_BadRequestAndSuccess(t *testing.T) {
	setupGinTest()

	mock := &mockRoutineService{}
	handler := NewRoutineHandler(mock)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("{invalid}"))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	handler.CreateRoutine(c)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for bad JSON, got %d", w.Code)
	}

	mock2 := &mockRoutineService{
		CreateRoutineFunc: func(ownerID string, input dto.RoutineRequest) (dto.RoutineResponse, error) {
			return dto.RoutineResponse{ID: "r1", UserID: ownerID, Name: input.Name, Excercises: input.Excercises}, nil
		},
	}
	handler2 := NewRoutineHandler(mock2)
	body := dto.RoutineRequest{Name: "name", Excercises: []dto.RoutineExcerciseList{{ExerciseID: "e1", Order: 1, Sets: 1, Reps: 1}}}
	bts, _ := json.Marshal(body)
	w2 := httptest.NewRecorder()
	c2, _ := gin.CreateTestContext(w2)
	req2 := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(bts))
	req2.Header.Set("Content-Type", "application/json")
	c2.Request = req2
	c2.Set("user_id", "owner1")
	handler2.CreateRoutine(c2)
	if w2.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d, body: %s", w2.Code, w2.Body.String())
	}
}

func TestUpdateAndDeleteRoutine_AuthorizationPaths(t *testing.T) {
	setupGinTest()

	mockUpd := &mockRoutineService{
		UpdateRoutineFunc: func(ownerID, routineID string, input dto.RoutineRequest) (dto.RoutineResponse, error) {
			return dto.RoutineResponse{}, errors.New("no autorizado: no es el owner de la rutina")
		},
	}
	hUpd := NewRoutineHandler(mockUpd)
	wUpd := httptest.NewRecorder()
	cUpd, _ := gin.CreateTestContext(wUpd)
	reqBody, _ := json.Marshal(dto.RoutineRequest{Name: "n", Excercises: []dto.RoutineExcerciseList{{ExerciseID: "e1", Order: 1, Sets: 1, Reps: 1}}})
	cUpd.Request = httptest.NewRequest(http.MethodPut, "/", bytes.NewBuffer(reqBody))
	cUpd.Request.Header.Set("Content-Type", "application/json")
	cUpd.Params = gin.Params{{Key: "id", Value: "r1"}}
	cUpd.Set("user_id", "owner1")
	hUpd.UpdateRoutine(cUpd)
	if wUpd.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for update forbidden, got %d, body: %s", wUpd.Code, wUpd.Body.String())
	}

	mockDel := &mockRoutineService{
		DeleteRoutineFunc: func(ownerID, routineID string) error { return nil },
	}
	hDel := NewRoutineHandler(mockDel)
	wDel := httptest.NewRecorder()
	cDel, _ := gin.CreateTestContext(wDel)
	cDel.Request = httptest.NewRequest(http.MethodDelete, "/", nil)
	cDel.Params = gin.Params{{Key: "id", Value: "r1"}}
	cDel.Set("user_id", "owner1")
	hDel.DeleteRoutine(cDel)
	if wDel.Code != http.StatusOK {
		t.Fatalf("expected 200 for delete success, got %d, body: %s", wDel.Code, wDel.Body.String())
	}

	mockDelForbidden := &mockRoutineService{
		DeleteRoutineFunc: func(ownerID, routineID string) error { return errors.New("no autorizado: no es el owner de la rutina") },
	}
	hDelF := NewRoutineHandler(mockDelForbidden)
	wDelF := httptest.NewRecorder()
	cDelF, _ := gin.CreateTestContext(wDelF)
	cDelF.Request = httptest.NewRequest(http.MethodDelete, "/", nil)
	cDelF.Params = gin.Params{{Key: "id", Value: "r1"}}
	cDelF.Set("user_id", "owner1")
	hDelF.DeleteRoutine(cDelF)
	if wDelF.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for delete forbidden, got %d, body: %s", wDelF.Code, wDelF.Body.String())
	}
}
