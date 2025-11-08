package repositories

import (
	"context"
	"backend/database"
	"backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LogRepositoryInterface interface {
	CreateLog(log models.Log) (*mongo.InsertOneResult, error)
	GetLogs(filter bson.M, limit, skip int64) ([]models.Log, error)
	GetLogByID(id string) (models.Log, error)
	DeleteAllLogs() (*mongo.DeleteResult, error)
}

type LogRepository struct {
	db database.DB
}

func NewLogRepository(db database.DB) *LogRepository {
	return &LogRepository{
		db: db,
	}
}

func (repository LogRepository) CreateLog(log models.Log) (*mongo.InsertOneResult, error) {
	collection := repository.db.GetClient().Database("fitness_db").Collection("logs")
	result, err := collection.InsertOne(context.TODO(), log)
	return result, err
}

func (repository LogRepository) GetLogs(filter bson.M, limit, skip int64) ([]models.Log, error) {
	collection := repository.db.GetClient().Database("fitness_db").Collection("logs")
	
	opts := options.Find()
	opts.SetSort(bson.D{{Key: "timestamp", Value: -1}})
	if limit > 0 {
		opts.SetLimit(limit)
	}
	if skip > 0 {
		opts.SetSkip(skip)
	}

	cursor, err := collection.Find(context.TODO(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var logs []models.Log
	for cursor.Next(context.Background()) {
		var log models.Log
		err := cursor.Decode(&log)
		if err != nil {
			continue
		}
		logs = append(logs, log)
	}
	return logs, nil
}

func (repository LogRepository) GetLogByID(id string) (models.Log, error) {
	collection := repository.db.GetClient().Database("fitness_db").Collection("logs")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Log{}, err
	}
	filter := bson.M{"_id": objectID}
	var log models.Log
	err = collection.FindOne(context.TODO(), filter).Decode(&log)
	return log, err
}

func (repository LogRepository) DeleteAllLogs() (*mongo.DeleteResult, error) {
	collection := repository.db.GetClient().Database("fitness_db").Collection("logs")
	filter := bson.M{}
	result, err := collection.DeleteMany(context.TODO(), filter)
	return result, err
}

