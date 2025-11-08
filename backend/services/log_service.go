package services

import (
	"backend/models"
	"backend/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type LogServiceInterface interface {
	CreateLog(level models.LogLevel, logType models.LogType, action string, userID *primitive.ObjectID, userName string, details string) error
	GetLogs(limit, skip int) ([]models.Log, error)
	GetLogByID(id string) (models.Log, error)
	DeleteAllLogs() error
}

type LogService struct {
	repo repositories.LogRepositoryInterface
}

func NewLogService(repo repositories.LogRepositoryInterface) *LogService {
	return &LogService{repo: repo}
}

func (s *LogService) CreateLog(level models.LogLevel, logType models.LogType, action string, userID *primitive.ObjectID, userName string, details string) error {
	now := time.Now()
	log := models.Log{
		ID:        primitive.NewObjectID(),
		Timestamp: now,
		Level:     level,
		Type:      logType,
		Action:    action,
		UserID:    userID,
		UserName:  userName,
		Details:   details,
		CreatedAt: now,
	}
	_, err := s.repo.CreateLog(log)
	return err
}

func (s *LogService) GetLogs(limit, skip int) ([]models.Log, error) {
	filter := bson.M{}

	var limit64 int64 = 0
	var skip64 int64 = 0
	if limit > 0 {
		limit64 = int64(limit)
	}
	if skip > 0 {
		skip64 = int64(skip)
	}

	return s.repo.GetLogs(filter, limit64, skip64)
}

func (s *LogService) GetLogByID(id string) (models.Log, error) {
	return s.repo.GetLogByID(id)
}

func (s *LogService) DeleteAllLogs() error {
	_, err := s.repo.DeleteAllLogs()
	return err
}

