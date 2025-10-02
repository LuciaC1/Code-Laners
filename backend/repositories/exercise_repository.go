package repositories

import (
	"context"

	"backend/database"
	"backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ExerciseRepositoryInterface interface {
	GetExercises(name, category, muscleGroup string) ([]models.Exercise, error)
	GetExerciseByID(id string) (models.Exercise, error)
	CreateExercise(exercise models.Exercise) (*mongo.InsertOneResult, error)
	UpdateExercise(exercise models.Exercise) (*mongo.UpdateResult, error)
	DeleteExercise(id primitive.ObjectID) (*mongo.DeleteResult, error)
}

type ExerciseRepository struct {
	db database.DB
}

func NewExerciseRepository(db database.DB) *ExerciseRepository {
	return &ExerciseRepository{
		db: db,
	}
}

func (repository ExerciseRepository) GetExercises(name, category, muscleGroup string) ([]models.Exercise, error) {
	collection := repository.db.GetClient().Database("fitness_db").Collection("exercises")

	filter := bson.M{}

	if name != "" {
		filter["name"] = bson.M{"$regex": name, "$options": "i"}
	}
	if category != "" {
		filter["category"] = category
	}
	if muscleGroup != "" {
		filter["muscle_group"] = muscleGroup
	}

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var exercises []models.Exercise
	for cursor.Next(context.Background()) {
		var exercise models.Exercise
		if err := cursor.Decode(&exercise); err != nil {
			continue
		}
		exercises = append(exercises, exercise)
	}

	return exercises, nil
}

func (repository ExerciseRepository) GetExerciseByID(id string) (models.Exercise, error) {
	collection := repository.db.GetClient().Database("fitness_db").Collection("exercises")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Exercise{}, err
	}

	filter := bson.M{"_id": objectID}
	var exercise models.Exercise

	err = collection.FindOne(context.TODO(), filter).Decode(&exercise)
	return exercise, err
}

func (repository ExerciseRepository) CreateExercise(exercise models.Exercise) (*mongo.InsertOneResult, error) {
	collection := repository.db.GetClient().Database("fitness_db").Collection("exercises")
	result, err := collection.InsertOne(context.TODO(), exercise)
	return result, err
}

func (repository ExerciseRepository) UpdateExercise(exercise models.Exercise) (*mongo.UpdateResult, error) {
	collection := repository.db.GetClient().Database("fitness_db").Collection("exercises")

	filter := bson.M{"_id": exercise.ID}
	update := bson.M{"$set": bson.M{
		"name":         exercise.Name,
		"description":  exercise.Description,
		"category":     exercise.Category,
		"muscle_group": exercise.MuscleGroup,
		"difficulty":   exercise.Difficulty,
		"media_url":    exercise.MediaURL,
		"steps":        exercise.Steps,
		"updated_at":   exercise.UpdatedAt,
	}}

	result, err := collection.UpdateOne(context.TODO(), filter, update)
	return result, err
}

func (repository ExerciseRepository) DeleteExercise(id primitive.ObjectID) (*mongo.DeleteResult, error) {
	collection := repository.db.GetClient().Database("fitness_db").Collection("exercises")

	filter := bson.M{"_id": id}
	result, err := collection.DeleteOne(context.TODO(), filter)
	return result, err
}
