package services

import (
	"testing"

	"backend/dto"
	"backend/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type mockUserRepo struct {
	getUserFn     func(name string) ([]models.User, error)
	getUserByIDFn func(id string) (models.User, error)
	createUserFn  func(user models.User) (*mongo.InsertOneResult, error)
	updateUserFn  func(user models.User) (*mongo.UpdateResult, error)
	deleteUserFn  func(id primitive.ObjectID) (*mongo.DeleteResult, error)
}

func (m *mockUserRepo) GetUser(name string) ([]models.User, error) {
	if m.getUserFn == nil {
		return nil, nil
	}
	return m.getUserFn(name)
}
func (m *mockUserRepo) GetUserByID(id string) (models.User, error) {
	if m.getUserByIDFn == nil {
		return models.User{}, nil
	}
	return m.getUserByIDFn(id)
}
func (m *mockUserRepo) CreateUser(user models.User) (*mongo.InsertOneResult, error) {
	if m.createUserFn == nil {
		return &mongo.InsertOneResult{InsertedID: user.ID}, nil
	}
	return m.createUserFn(user)
}
func (m *mockUserRepo) UpdateUser(user models.User) (*mongo.UpdateResult, error) {
	if m.updateUserFn == nil {
		return &mongo.UpdateResult{}, nil
	}
	return m.updateUserFn(user)
}
func (m *mockUserRepo) DeleteUser(id primitive.ObjectID) (*mongo.DeleteResult, error) {
	if m.deleteUserFn == nil {
		return &mongo.DeleteResult{}, nil
	}
	return m.deleteUserFn(id)
}

func TestRegister_Success(t *testing.T) {
	repo := &mockUserRepo{
		getUserFn: func(name string) ([]models.User, error) { return []models.User{}, nil },
	}

	svc := NewUserService(repo)

	req := dto.RegisterRequest{
		Name:        "Alice",
		Email:       "alice@example.com",
		Password:    "secret",
		DateOfBirth: "2000-01-01",
		Level:       "beginner",
		Goals:       []string{"fitness"},
	}

	u, err := svc.Register(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if u.Email != req.Email || u.Name != req.Name {
		t.Fatalf("unexpected user returned: %+v", u)
	}
}

func TestRegister_DuplicateEmail(t *testing.T) {
	existing := models.User{ID: primitive.NewObjectID(), Email: "bob@example.com"}
	repo := &mockUserRepo{
		getUserFn: func(name string) ([]models.User, error) { return []models.User{existing}, nil },
	}
	svc := NewUserService(repo)
	req := dto.RegisterRequest{Name: "Bob", Email: "bob@example.com", Password: "p", DateOfBirth: "2000-01-01"}
	_, err := svc.Register(req)
	if err == nil {
		t.Fatalf("expected duplicate email error")
	}
}

func TestLogin_Success(t *testing.T) {
	pw := "mypw"
	hash, _ := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	stored := models.User{ID: primitive.NewObjectID(), Email: "c@example.com", PasswordHash: string(hash)}
	repo := &mockUserRepo{getUserFn: func(name string) ([]models.User, error) { return []models.User{stored}, nil }}
	svc := NewUserService(repo)
	got, err := svc.Login(dto.LoginRequest{Email: "c@example.com", Password: pw})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Email != stored.Email {
		t.Fatalf("expected email %s got %s", stored.Email, got.Email)
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("right"), bcrypt.DefaultCost)
	stored := models.User{ID: primitive.NewObjectID(), Email: "d@example.com", PasswordHash: string(hash)}
	repo := &mockUserRepo{getUserFn: func(name string) ([]models.User, error) { return []models.User{stored}, nil }}
	svc := NewUserService(repo)
	_, err := svc.Login(dto.LoginRequest{Email: "d@example.com", Password: "wrong"})
	if err == nil {
		t.Fatalf("expected wrong password error")
	}
}

func TestGetUsers_Mapping(t *testing.T) {
	users := []models.User{{ID: primitive.NewObjectID(), Name: "U1"}, {ID: primitive.NewObjectID(), Name: "U2"}}
	repo := &mockUserRepo{getUserFn: func(name string) ([]models.User, error) { return users, nil }}
	svc := NewUserService(repo)
	out, err := svc.GetUsers("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != len(users) {
		t.Fatalf("expected %d users got %d", len(users), len(out))
	}
}

func TestChangePassword_WrongOldPassword(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("oldpass"), bcrypt.DefaultCost)
	m := models.User{ID: primitive.NewObjectID(), PasswordHash: string(hash)}
	repo := &mockUserRepo{getUserByIDFn: func(id string) (models.User, error) { return m, nil }}
	svc := NewUserService(repo)
	err := svc.ChangePassword(m.ID.Hex(), dto.ChangePasswordRequest{OldPassword: "bad", NewPassword: "new"})
	if err == nil {
		t.Fatalf("expected error for wrong old password")
	}
}

func TestChangePassword_Success(t *testing.T) {
	old := "oldpass"
	hash, _ := bcrypt.GenerateFromPassword([]byte(old), bcrypt.DefaultCost)
	m := models.User{ID: primitive.NewObjectID(), PasswordHash: string(hash)}
	updated := false
	repo := &mockUserRepo{
		getUserByIDFn: func(id string) (models.User, error) { return m, nil },
		updateUserFn:  func(user models.User) (*mongo.UpdateResult, error) { updated = true; return &mongo.UpdateResult{}, nil },
	}
	svc := NewUserService(repo)
	err := svc.ChangePassword(m.ID.Hex(), dto.ChangePasswordRequest{OldPassword: old, NewPassword: "brandnew"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !updated {
		t.Fatalf("expected update to be called")
	}
}

func TestDeleteUser_InvalidHex(t *testing.T) {
	svc := NewUserService(&mockUserRepo{})
	err := svc.DeleteUser("nothex")
	if err == nil {
		t.Fatalf("expected error for invalid hex id")
	}
}
