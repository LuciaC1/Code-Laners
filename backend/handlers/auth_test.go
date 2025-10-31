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

	"golang.org/x/crypto/bcrypt"
)

type mockUserService struct {
	registerFn func(req dto.RegisterRequest) (string, error)
	loginFn    func(req dto.LoginRequest) (dto.User, error)
}

func (m *mockUserService) Register(req dto.RegisterRequest) (string, error) {
	if m.registerFn != nil {
		return m.registerFn(req)
	}
	return "", nil
}
func (m *mockUserService) Login(req dto.LoginRequest) (dto.User, error) {
	if m.loginFn != nil {
		return m.loginFn(req)
	}
	return dto.User{}, nil
}

func (m *mockUserService) GetUsers(name string) ([]dto.User, error)                      { return nil, nil }
func (m *mockUserService) GetUserByID(id string) (dto.User, error)                       { return dto.User{}, nil }
func (m *mockUserService) UpdateUser(id string, req dto.UpdateUserRequest) error         { return nil }
func (m *mockUserService) ChangePassword(id string, req dto.ChangePasswordRequest) error { return nil }
func (m *mockUserService) DeleteUser(id string) error                                    { return nil }

func makeReq(t *testing.T, method, path string, body interface{}) (*gin.Context, *httptest.ResponseRecorder) {
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			t.Fatalf("encode body: %v", err)
		}
	}
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	return c, w
}

func TestRegister_Success(t *testing.T) {

	gin.SetMode(gin.TestMode)

	now := time.Now()
	mockedUser := dto.User{
		ID:          primitive.NewObjectID(),
		Name:        "Alice",
		Email:       "alice@example.com",
		Role:        string(models.RoleUser),
		DateOfBirth: now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	msvc := &mockUserService{
		registerFn: func(req dto.RegisterRequest) (string, error) {

			return mockedUser.Email, nil
		},
	}

	handler := NewUserHandler(msvc)

	reqBody := dto.RegisterRequest{
		Name:        "Alice",
		Email:       "alice@example.com",
		Password:    "secret123",
		DateOfBirth: now.Format(time.RFC3339),
	}

	c, w := makeReq(t, "POST", "/register", reqBody)

	handler.Register(c)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status %d got %d", http.StatusCreated, w.Code)
	}

	var got string
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if got != mockedUser.Email {
		t.Fatalf("unexpected response: %v", got)
	}
}

func TestLogin_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	password := "password123"
	pwHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	userID := primitive.NewObjectID()
	now := time.Now()
	returnedUser := dto.User{
		ID:           userID,
		Name:         "Bob",
		Email:        "bob@example.com",
		PasswordHash: string(pwHash),
		Role:         string(models.RoleUser),
		DateOfBirth:  now,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	msvc := &mockUserService{
		loginFn: func(req dto.LoginRequest) (dto.User, error) {

			if req.Email == returnedUser.Email {
				return returnedUser, nil
			}
			return dto.User{}, nil
		},
	}

	handler := NewUserHandler(msvc)

	reqBody := dto.LoginRequest{
		Email:    returnedUser.Email,
		Password: password,
	}

	c, w := makeReq(t, "POST", "/login", reqBody)

	handler.Login(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d got %d body: %s", http.StatusOK, w.Code, w.Body.String())
	}
	var resp dto.AuthResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if resp.User.Email != returnedUser.Email {
		t.Fatalf("unexpected user in response: %+v", resp.User)
	}
}
