package repositories

import (
	"backend/database"
	"backend/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type StatsRepository struct {
	db database.DB
}

func NewStatsRepository(db database.DB) *StatsRepository {
	return &StatsRepository{db: db}
}

// Count documents helpers
func (r *StatsRepository) CountWorkoutsByUser(userID primitive.ObjectID) (int64, error) {
	collection := r.db.GetClient().Database("fitness_db").Collection("workouts")
	filter := bson.M{"user_id": userID}
	return collection.CountDocuments(context.TODO(), filter)
}

func (r *StatsRepository) SumWorkoutFieldsByUser(userID primitive.ObjectID) (int64, int64, error) {
	collection := r.db.GetClient().Database("fitness_db").Collection("workouts")
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.M{"user_id": userID}}},
		bson.D{{Key: "$group", Value: bson.M{"_id": nil, "totalMinutes": bson.M{"$sum": "$duration_minutes"}, "totalCalories": bson.M{"$sum": "$estimated_calories"}}}},
	}
	cursor, err := collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return 0, 0, err
	}
	defer cursor.Close(context.TODO())
	var res []bson.M
	if err := cursor.All(context.TODO(), &res); err != nil {
		return 0, 0, err
	}
	if len(res) == 0 {
		return 0, 0, nil
	}
	totalMinutes, _ := res[0]["totalMinutes"].(int32)
	totalCalories, _ := res[0]["totalCalories"].(int32)
	return int64(totalMinutes), int64(totalCalories), nil
}

func (r *StatsRepository) AggregateWorkoutsByWeek(userID primitive.ObjectID) ([]bson.M, error) {
	collection := r.db.GetClient().Database("fitness_db").Collection("workouts")
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.M{"user_id": userID}}},
		bson.D{{Key: "$addFields", Value: bson.M{"year": bson.M{"$year": "$completed_at"}, "week": bson.M{"$isoWeek": "$completed_at"}}}},
		bson.D{{Key: "$group", Value: bson.M{"_id": bson.M{"year": "$year", "week": "$week"}, "count": bson.M{"$sum": 1}}}},
		bson.D{{Key: "$sort", Value: bson.M{"_id.year": 1, "_id.week": 1}}},
	}
	cursor, err := collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())
	var out []bson.M
	if err := cursor.All(context.TODO(), &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *StatsRepository) AggregateWorkoutsByRoutine(userID primitive.ObjectID) ([]bson.M, error) {
	collection := r.db.GetClient().Database("fitness_db").Collection("workouts")
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.M{"user_id": userID}}},
		bson.D{{Key: "$group", Value: bson.M{"_id": "$routine_id", "count": bson.M{"$sum": 1}}}},
		bson.D{{Key: "$lookup", Value: bson.M{"from": "routines", "localField": "_id", "foreignField": "_id", "as": "routine"}}},
		bson.D{{Key: "$unwind", Value: bson.M{"path": "$routine", "preserveNullAndEmptyArrays": true}}},
		bson.D{{Key: "$project", Value: bson.M{"routine_id": "$_id", "routine_name": bson.M{"$ifNull": []interface{}{"$routine.name", "Sin rutina"}}, "count": 1}}},
		bson.D{{Key: "$sort", Value: bson.M{"count": -1}}},
	}
	cursor, err := collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())
	var out []bson.M
	if err := cursor.All(context.TODO(), &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *StatsRepository) GetProgressPoints(userID primitive.ObjectID, limit int) ([]models.Workout, error) {
	collection := r.db.GetClient().Database("fitness_db").Collection("workouts")
	filter := bson.M{"user_id": userID}
	opts := options.Find()
	// sort ascending by date
	sort := bson.D{{Key: "completed_at", Value: 1}}
	opts.SetSort(sort)
	if limit > 0 {
		opts.SetLimit(int64(limit))
	}
	cursor, err := collection.Find(context.TODO(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())
	var out []models.Workout
	for cursor.Next(context.TODO()) {
		var w models.Workout
		if err := cursor.Decode(&w); err != nil {
			continue
		}
		out = append(out, w)
	}
	return out, nil
}

func (r *StatsRepository) GetRecentWorkouts(userID primitive.ObjectID, limit int) ([]models.Workout, error) {
	collection := r.db.GetClient().Database("fitness_db").Collection("workouts")
	filter := bson.M{"user_id": userID}
	opts := options.Find()
	sort := bson.D{{Key: "completed_at", Value: -1}}
	opts.SetSort(sort)
	if limit > 0 {
		opts.SetLimit(int64(limit))
	}
	cursor, err := collection.Find(context.TODO(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())
	var out []models.Workout
	for cursor.Next(context.TODO()) {
		var w models.Workout
		if err := cursor.Decode(&w); err != nil {
			continue
		}
		out = append(out, w)
	}
	return out, nil
}

// Admin helpers
func (r *StatsRepository) CountCollection(name string) (int64, error) {
	collection := r.db.GetClient().Database("fitness_db").Collection(name)
	return collection.CountDocuments(context.TODO(), bson.M{})
}

func (r *StatsRepository) UsersByMonth() ([]bson.M, error) {
	collection := r.db.GetClient().Database("fitness_db").Collection("users")
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$addFields", Value: bson.M{"year": bson.M{"$year": "$created_at"}, "month": bson.M{"$month": "$created_at"}}}},
		bson.D{{Key: "$group", Value: bson.M{"_id": bson.M{"year": "$year", "month": "$month"}, "count": bson.M{"$sum": 1}}}},
		bson.D{{Key: "$sort", Value: bson.M{"_id.year": -1, "_id.month": -1}}},
	}
	cursor, err := collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())
	var out []bson.M
	if err := cursor.All(context.TODO(), &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *StatsRepository) ExercisesByCategory() ([]bson.M, error) {
	collection := r.db.GetClient().Database("fitness_db").Collection("exercises")
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$group", Value: bson.M{"_id": "$category", "count": bson.M{"$sum": 1}}}},
		bson.D{{Key: "$sort", Value: bson.M{"count": -1}}},
	}
	cursor, err := collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())
	var out []bson.M
	if err := cursor.All(context.TODO(), &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *StatsRepository) PopularRoutines(limit int) ([]bson.M, error) {
	collection := r.db.GetClient().Database("fitness_db").Collection("workouts")
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$group", Value: bson.M{"_id": "$routine_id", "usage": bson.M{"$sum": 1}}}},
		bson.D{{Key: "$lookup", Value: bson.M{"from": "routines", "localField": "_id", "foreignField": "_id", "as": "routine"}}},
		bson.D{{Key: "$unwind", Value: bson.M{"path": "$routine", "preserveNullAndEmptyArrays": true}}},
		bson.D{{Key: "$project", Value: bson.M{"routine_id": "$_id", "name": "$routine.name", "usage": 1}}},
		bson.D{{Key: "$sort", Value: bson.M{"usage": -1}}},
	}
	if limit > 0 {
		pipeline = append(pipeline, bson.D{{Key: "$limit", Value: limit}})
	}
	cursor, err := collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())
	var out []bson.M
	if err := cursor.All(context.TODO(), &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *StatsRepository) RecentWorkoutsGlobal(limit int) ([]models.Workout, error) {
	collection := r.db.GetClient().Database("fitness_db").Collection("workouts")
	opts := options.Find()
	sort := bson.D{{Key: "completed_at", Value: -1}}
	opts.SetSort(sort)
	if limit > 0 {
		opts.SetLimit(int64(limit))
	}
	cursor, err := collection.Find(context.TODO(), bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())
	var out []models.Workout
	for cursor.Next(context.TODO()) {
		var w models.Workout
		if err := cursor.Decode(&w); err != nil {
			continue
		}
		out = append(out, w)
	}
	return out, nil
}
