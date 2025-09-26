package repositories

import (
	"context"

	"backend/database"
	"backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type WorkoutRepositoryInterface interface {
	GetWorkouts(userID primitive.ObjectID) ([]models.Workout, error)
	GetWorkoutByID(id string) (models.Workout, error)
	CreateWorkout(workout models.Workout) (*mongo.InsertOneResult, error)
	UpdateWorkout(workout models.Workout) (*mongo.UpdateResult, error)
	DeleteWorkout(id primitive.ObjectID) (*mongo.DeleteResult, error)
}

type WorkoutRepository struct {
	db database.DB
}

func NewWorkoutRepository(db database.DB) *WorkoutRepository {
	return &WorkoutRepository{
		db: db,
	}
}

func (repository WorkoutRepository) GetWorkouts(userID primitive.ObjectID) ([]models.Workout, error) {
	collection := repository.db.GetClient().Database("fitness_db").Collection("workouts")

	filter := bson.M{"user_id": userID}

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var workouts []models.Workout
	for cursor.Next(context.Background()) {
		var workout models.Workout
		if err := cursor.Decode(&workout); err != nil {
			continue
		}
		workouts = append(workouts, workout)
	}

	return workouts, nil
}

func (repository WorkoutRepository) GetWorkoutByID(id string) (models.Workout, error) {
	collection := repository.db.GetClient().Database("fitness_db").Collection("workouts")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Workout{}, err
	}

	filter := bson.M{"_id": objectID}
	var workout models.Workout

	err = collection.FindOne(context.TODO(), filter).Decode(&workout)
	return workout, err
}

func (repository WorkoutRepository) CreateWorkout(workout models.Workout) (*mongo.InsertOneResult, error) {
	collection := repository.db.GetClient().Database("fitness_db").Collection("workouts")
	result, err := collection.InsertOne(context.TODO(), workout)
	return result, err
}

func (repository WorkoutRepository) UpdateWorkout(workout models.Workout) (*mongo.UpdateResult, error) {
	collection := repository.db.GetClient().Database("fitness_db").Collection("workouts")

	filter := bson.M{"_id": workout.ID}
	update := bson.M{"$set": bson.M{
		"user_id":             workout.UserID,
		"routine_id":          workout.RoutineID,
		"performed_exercises": workout.PerformedExercises,
		"completed_at":        workout.CompletedAt,
		"duration_minutes":    workout.DurationMinutes,
		"notes":               workout.Notes,
	}}

	result, err := collection.UpdateOne(context.TODO(), filter, update)
	return result, err
}

func (repository WorkoutRepository) DeleteWorkout(id primitive.ObjectID) (*mongo.DeleteResult, error) {
	collection := repository.db.GetClient().Database("fitness_db").Collection("workouts")

	filter := bson.M{"_id": id}
	result, err := collection.DeleteOne(context.TODO(), filter)
	return result, err
}
