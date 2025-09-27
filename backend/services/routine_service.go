package services

import (
	"errors"
	"time"

	"backend/models"
	"backend/repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RoutineService struct {
	repo repositories.RoutineRepositoryInterface
}

func NewRoutineService(repo repositories.RoutineRepositoryInterface) *RoutineService {
	return &RoutineService{repo: repo}
}

func (s *RoutineService) GetRoutines(ownerHex, name string) ([]models.Routine, error) {
	if ownerHex == "" {
		return nil, errors.New("owner id required")
	}
	ownerID, err := primitive.ObjectIDFromHex(ownerHex)
	if err != nil {
		return nil, err
	}
	return s.repo.GetRoutines(ownerID, name)
}

func (s *RoutineService) GetRoutineByID(id string) (models.Routine, error) {
	if id == "" {
		return models.Routine{}, errors.New("id required")
	}
	return s.repo.GetRoutineByID(id)
}

func (s *RoutineService) CreateRoutine(r models.Routine, ownerHex string) (*mongo.InsertOneResult, error) {
	if ownerHex == "" {
		return nil, errors.New("owner id required")
	}
	if r.Name == "" {
		return nil, errors.New("routine name is required")
	}
	ownerID, err := primitive.ObjectIDFromHex(ownerHex)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	r.OwnerID = ownerID
	r.CreatedAt = now
	r.UpdatedAt = now
	if r.Entries == nil {
		r.Entries = []models.RoutineEntry{}
	}
	return s.repo.CreateRoutine(r)
}

func (s *RoutineService) UpdateRoutine(idHex string, payload models.Routine, ownerHex string) (*mongo.UpdateResult, error) {
	if idHex == "" {
		return nil, errors.New("id required")
	}
	if ownerHex == "" {
		return nil, errors.New("owner id required")
	}
	existing, err := s.repo.GetRoutineByID(idHex)
	if err != nil {
		return nil, err
	}
	ownerID, err := primitive.ObjectIDFromHex(ownerHex)
	if err != nil {
		return nil, err
	}
	if existing.OwnerID != ownerID {
		return nil, errors.New("forbidden: not the owner")
	}
	if payload.Name != "" {
		existing.Name = payload.Name
	}
	if payload.Description != "" {
		existing.Description = payload.Description
	}
	if payload.IsPublic {
		existing.IsPublic = payload.IsPublic
	}
	if len(payload.Entries) > 0 {
		existing.Entries = payload.Entries
	}
	existing.UpdatedAt = time.Now()
	return s.repo.UpdateRoutine(existing)
}

func (s *RoutineService) DeleteRoutine(idHex, ownerHex string) (*mongo.DeleteResult, error) {
	if idHex == "" {
		return nil, errors.New("id required")
	}
	if ownerHex == "" {
		return nil, errors.New("owner id required")
	}
	existing, err := s.repo.GetRoutineByID(idHex)
	if err != nil {
		return nil, err
	}
	ownerID, err := primitive.ObjectIDFromHex(ownerHex)
	if err != nil {
		return nil, err
	}
	if existing.OwnerID != ownerID {
		return nil, errors.New("forbidden: not the owner")
	}
	oid, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return nil, err
	}
	return s.repo.DeleteRoutine(oid)
}

func (s *RoutineService) DuplicateRoutine(idHex, ownerHex, newName string) (*mongo.InsertOneResult, error) {
	if idHex == "" {
		return nil, errors.New("id required")
	}
	if ownerHex == "" {
		return nil, errors.New("owner id required")
	}
	orig, err := s.repo.GetRoutineByID(idHex)
	if err != nil {
		return nil, err
	}
	ownerID, err := primitive.ObjectIDFromHex(ownerHex)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	copyRoutine := models.Routine{
		OwnerID:     ownerID,
		Name:        newNameOrDefault(newName, orig.Name),
		Description: orig.Description,
		Entries:     duplicateEntries(orig.Entries),
		IsPublic:    orig.IsPublic,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	return s.repo.CreateRoutine(copyRoutine)
}

func (s *RoutineService) AddEntry(routineIDHex string, entry models.RoutineEntry, ownerHex string) (*mongo.UpdateResult, error) {
	if routineIDHex == "" {
		return nil, errors.New("routine id required")
	}
	if ownerHex == "" {
		return nil, errors.New("owner id required")
	}
	routine, err := s.repo.GetRoutineByID(routineIDHex)
	if err != nil {
		return nil, err
	}
	ownerID, err := primitive.ObjectIDFromHex(ownerHex)
	if err != nil {
		return nil, err
	}
	if routine.OwnerID != ownerID {
		return nil, errors.New("forbidden: not the owner")
	}
	routine.Entries = append(routine.Entries, entry)
	routine.UpdatedAt = time.Now()
	return s.repo.UpdateRoutine(routine)
}

func (s *RoutineService) UpdateEntry(routineIDHex string, index int, entry models.RoutineEntry, ownerHex string) (*mongo.UpdateResult, error) {
	if routineIDHex == "" {
		return nil, errors.New("routine id required")
	}
	if ownerHex == "" {
		return nil, errors.New("owner id required")
	}
	routine, err := s.repo.GetRoutineByID(routineIDHex)
	if err != nil {
		return nil, err
	}
	ownerID, err := primitive.ObjectIDFromHex(ownerHex)
	if err != nil {
		return nil, err
	}
	if routine.OwnerID != ownerID {
		return nil, errors.New("forbidden: not the owner")
	}
	if index < 0 || index >= len(routine.Entries) {
		return nil, errors.New("entry index out of range")
	}
	routine.Entries[index] = entry
	routine.UpdatedAt = time.Now()
	return s.repo.UpdateRoutine(routine)
}

func (s *RoutineService) RemoveEntry(routineIDHex string, index int, ownerHex string) (*mongo.UpdateResult, error) {
	if routineIDHex == "" {
		return nil, errors.New("routine id required")
	}
	if ownerHex == "" {
		return nil, errors.New("owner id required")
	}
	routine, err := s.repo.GetRoutineByID(routineIDHex)
	if err != nil {
		return nil, err
	}
	ownerID, err := primitive.ObjectIDFromHex(ownerHex)
	if err != nil {
		return nil, err
	}
	if routine.OwnerID != ownerID {
		return nil, errors.New("forbidden: not the owner")
	}
	if index < 0 || index >= len(routine.Entries) {
		return nil, errors.New("entry index out of range")
	}
	routine.Entries = append(routine.Entries[:index], routine.Entries[index+1:]...)
	routine.UpdatedAt = time.Now()
	return s.repo.UpdateRoutine(routine)
}

func newNameOrDefault(newName, original string) string {
	if newName != "" {
		return newName
	}
	return "Copia de " + original
}

func duplicateEntries(orig []models.RoutineEntry) []models.RoutineEntry {
	if len(orig) == 0 {
		return []models.RoutineEntry{}
	}
	cp := make([]models.RoutineEntry, len(orig))
	copy(cp, orig)
	return cp
}
