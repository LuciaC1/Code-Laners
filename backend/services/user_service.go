package services

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"

	"backend/models"
	"backend/repositories"
)

type UserService struct {
	Repo *repositories.UserRepository
}

func NewUserService(repo *repositories.UserRepository) *UserService {
	return &UserService{Repo: repo}
}

func (s *UserService) GetUsers(name string) ([]models.User, error) {
	return s.Repo.GetUser(name)
}

func (s *UserService) GetUserByID(id string) (models.User, error) {
	return s.Repo.GetUserByID(id)
}

func (s *UserService) RegisterUser(u models.User, plainPassword string) (*primitive.ObjectID, error) {
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

	now := time.Now()
	u.CreatedAt = now
	u.UpdatedAt = now

	res, err := s.Repo.CreateUser(u)
	if err != nil {
		return nil, err
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, errors.New("could not parse inserted id")
	}
	return &oid, nil
}

func (s *UserService) UpdateProfile(idHex string, payload models.User) error {
	existing, err := s.Repo.GetUserByID(idHex)
	if err != nil {
		return err
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
		existing.Level = payload.Level
	}
	if payload.Goals != nil && len(payload.Goals) > 0 {
		existing.Goals = payload.Goals
	}
	existing.UpdatedAt = time.Now()

	_, err = s.Repo.UpdaterUser(existing)
	return err
}

func (s *UserService) ChangePassword(idHex string, oldPassword, newPassword string) error {
	if len(newPassword) < 6 {
		return errors.New("new password must be at least 6 characters")
	}

	user, err := s.Repo.GetUserByID(idHex)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword)); err != nil {
		return errors.New("old password is incorrect")
	}

	newHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(newHash)
	user.UpdatedAt = time.Now()

	_, err = s.Repo.UpdaterUser(user)
	return err
}

func (s *UserService) DeleteUserByHex(idHex string) error {
	oid, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return err
	}
	_, err = s.Repo.DeleteUser(oid)
	return err
}
