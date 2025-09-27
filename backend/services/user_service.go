package services

import (
	"errors"
	"time"

	"backend/models"
	"backend/repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repositories.UserRepository
}

func NewUserService(repo *repositories.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUsers(name string) ([]models.User, error) {
	return s.repo.GetUser(name)
}

func (s *UserService) GetUserByID(id string) (models.User, error) {
	if id == "" {
		return models.User{}, errors.New("id required")
	}
	return s.repo.GetUserByID(id)
}

func (s *UserService) CreateUser(u models.User, plainPassword string) (*mongo.InsertOneResult, error) {
	if u.Name == "" || u.Email == "" {
		return nil, errors.New("name and email are required")
	}
	if len(plainPassword) < 6 {
		return nil, errors.New("password must be at least 6 characters")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	u.PasswordHash = string(hash)
	if u.Role == "" {
		u.Role = "user"
	}
	now := time.Now()
	u.CreatedAt = now
	u.UpdatedAt = now
	return s.repo.CreateUser(u)
}

func normalizeLevel(l string) string {
	switch l {
	case "principiante":
		return "beginner"
	case "intermedio":
		return "intermediate"
	case "avanzado":
		return "advanced"
	default:
		return l
	}
}

func (s *UserService) UpdateProfile(id string, payload models.User) (*mongo.UpdateResult, error) {
	if id == "" {
		return nil, errors.New("id required")
	}
	existing, err := s.repo.GetUserByID(id)
	if err != nil {
		return nil, err
	}
	if payload.Name != "" {
		existing.Name = payload.Name
	}
	if payload.Email != "" {
		existing.Email = payload.Email
	}
	if !payload.DateOfBirth.IsZero() {
		existing.DateOfBirth = payload.DateOfBirth
	}
	if payload.Weight != 0 {
		existing.Weight = payload.Weight
	}
	if payload.Height != 0 {
		existing.Height = payload.Height
	}
	if payload.Level != "" {
		existing.Level = normalizeLevel(payload.Level)
	}
	if len(payload.Goals) > 0 {
		existing.Goals = payload.Goals
	}
	existing.UpdatedAt = time.Now()
	return s.repo.UpdaterUser(existing)
}

func (s *UserService) ChangePassword(id, oldPassword, newPassword string) (*mongo.UpdateResult, error) {
	if id == "" {
		return nil, errors.New("id required")
	}
	if len(newPassword) < 6 {
		return nil, errors.New("new password must be at least 6 characters")
	}
	user, err := s.repo.GetUserByID(id)
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword)); err != nil {
		return nil, errors.New("old password is incorrect")
	}
	newHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.PasswordHash = string(newHash)
	user.UpdatedAt = time.Now()
	return s.repo.UpdaterUser(user)
}

func (s *UserService) DeleteUser(id string) (*mongo.DeleteResult, error) {
	if id == "" {
		return nil, errors.New("id required")
	}
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	return s.repo.DeleteUser(oid)
}
