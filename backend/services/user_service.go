package services

import (
	"errors"
	"regexp"
	"time"

	"backend/dto"
	"backend/models"
	"backend/repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceInterface interface {
	Register(req dto.RegisterRequest) (string, error)
	Login(req dto.LoginRequest) (dto.User, error)
	GetUsers(name string) ([]dto.User, error)
	GetUserByID(id string) (dto.User, error)
	UpdateUser(id string, req dto.UpdateUserRequest) error
	ChangePassword(id string, req dto.ChangePasswordRequest) error
	DeleteUser(id string) error
}

type UserService struct {
	repo repositories.UserRepositoryInterface
}

func NewUserService(repo repositories.UserRepositoryInterface) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(req dto.RegisterRequest) (dto.User, error) {
	if req.Name == "" || req.Email == "" || req.Password == "" || req.DateOfBirth == "" {
		return dto.User{}, errors.New("datos incompletos")
	}
	if !isValidEmail(req.Email) {
		return dto.User{}, errors.New("email inválido")
	}
	dob, err := time.Parse(time.RFC3339, req.DateOfBirth)
	if err != nil {
		dob, err = time.Parse("2006-01-02", req.DateOfBirth)
		if err != nil {
			return dto.User{}, errors.New("date_of_birth formato inválido, use ISO")
		}
	}
	candidates, err := s.repo.GetUser("")
	if err == nil {
		for _, u := range candidates {
			if u.Email == req.Email {
				return dto.User{}, errors.New("email ya registrado")
			}
		}
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return dto.User{}, err
	}
	now := time.Now()
	user := models.User{
		ID:           primitive.NewObjectID(),
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hash),
		Role:         "user",
		DateOfBirth:  dob,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if req.Weight != nil {
		user.Weight = *req.Weight
	}
	if req.Height != nil {
		user.Height = *req.Height
	}
	user.Level = req.Level
	user.Goals = req.Goals

	res, err := s.repo.CreateUser(user)
	if err != nil {
		return dto.User{}, err
	}
	if res.InsertedID == nil {
		return dto.User{}, errors.New("no se pudo crear el usuario")
	}
	return modelUserToDTO(user), nil

}
func (s *UserService) Login(req dto.LoginRequest) (dto.User, error) {
	if req.Email == "" || req.Password == "" {
		return dto.User{}, errors.New("credenciales requeridas")
	}
	candidates, err := s.repo.GetUser("")
	if err != nil {
		return dto.User{}, err
	}
	var found *models.User
	for _, u := range candidates {
		if u.Email == req.Email {
			found = &u
			break
		}
	}
	if found == nil {
		return dto.User{}, errors.New("usuario no encontrado")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(found.PasswordHash), []byte(req.Password)); err != nil {
		return dto.User{}, errors.New("contraseña incorrecta")
	}
	return modelUserToDTO(*found), nil
}

func (s *UserService) GetUsers(name string) ([]dto.User, error) {
	models, err := s.repo.GetUser(name)
	if err != nil {
		return nil, err
	}
	out := make([]dto.User, 0, len(models))
	for _, m := range models {
		out = append(out, modelUserToDTO(m))
	}
	return out, nil
}

func (s *UserService) GetUserByID(id string) (dto.User, error) {
	m, err := s.repo.GetUserByID(id)
	if err != nil {
		return dto.User{}, err
	}
	return modelUserToDTO(m), nil
}

func (s *UserService) UpdateUser(id string, req dto.UpdateUserRequest) error {
	m, err := s.repo.GetUserByID(id)
	if err != nil {
		return err
	}
	if req.Email != "" && !isValidEmail(req.Email) {
		return errors.New("email inválido")
	}
	if req.Name != "" {
		m.Name = req.Name
	}
	if req.Email != "" {
		m.Email = req.Email
	}
	if req.Weight != 0 {
		m.Weight = req.Weight
	}
	if req.Height != 0 {
		m.Height = req.Height
	}
	if req.Level != "" {
		m.Level = req.Level
	}
	if req.Goals != nil {
		m.Goals = req.Goals
	}
	m.UpdatedAt = time.Now()
	_, err = s.repo.UpdateUser(m)
	return err
}

func (s *UserService) ChangePassword(id string, req dto.ChangePasswordRequest) error {
	m, err := s.repo.GetUserByID(id)
	if err != nil {
		return err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(m.PasswordHash), []byte(req.OldPassword)); err != nil {
		return errors.New("contraseña actual incorrecta")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	m.PasswordHash = string(hash)
	m.UpdatedAt = time.Now()
	_, err = s.repo.UpdateUser(m)
	return err
}

func (s *UserService) DeleteUser(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = s.repo.DeleteUser(objID)
	return err
}

func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func modelUserToDTO(m models.User) dto.User {
	return dto.User{
		ID:          m.ID,
		Name:        m.Name,
		Email:       m.Email,
		Role:        string(m.Role),
		DateOfBirth: m.DateOfBirth,
		Weight:      m.Weight,
		Height:      m.Height,
		Level:       m.Level,
		Goals:       m.Goals,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}
