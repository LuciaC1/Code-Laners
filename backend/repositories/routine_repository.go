package repositories

import (
	"context"

	"backend/database"
	"backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RoutineRepositoryInterface interface {
	GetRoutines(ownerID primitive.ObjectID, name string) ([]models.Routine, error)
	GetRoutineByID(id string) (models.Routine, error)
	CreateRoutine(routine models.Routine) (*mongo.InsertOneResult, error)
	UpdateRoutine(routine models.Routine) (*mongo.UpdateResult, error)
	DeleteRoutine(id primitive.ObjectID) (*mongo.DeleteResult, error)
}

type RoutineRepository struct {
	db database.DB
}

func NewRoutineRepository(db database.DB) *RoutineRepository {
	return &RoutineRepository{
		db: db,
	}
}

func (repository RoutineRepository) GetRoutines(ownerID primitive.ObjectID, name string) ([]models.Routine, error) {
	collection := repository.db.GetClient().Database("fitness_db").Collection("routines")

	filter := bson.M{"owner_id": ownerID}
	if name != "" {
		filter["name"] = bson.M{"$regex": name, "$options": "i"}
	}

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var routines []models.Routine
	for cursor.Next(context.Background()) {
		var routine models.Routine
		if err := cursor.Decode(&routine); err != nil {
			continue
		}
		routines = append(routines, routine)
	}

	return routines, nil
}

func (repository RoutineRepository) GetRoutineByID(id string) (models.Routine, error) {
	collection := repository.db.GetClient().Database("fitness_db").Collection("routines")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Routine{}, err
	}

	filter := bson.M{"_id": objectID}
	var routine models.Routine

	err = collection.FindOne(context.TODO(), filter).Decode(&routine)
	return routine, err
}

func (repository RoutineRepository) CreateRoutine(routine models.Routine) (*mongo.InsertOneResult, error) {
	collection := repository.db.GetClient().Database("fitness_db").Collection("routines")
	result, err := collection.InsertOne(context.TODO(), routine)
	return result, err
}

func (repository RoutineRepository) UpdateRoutine(routine models.Routine) (*mongo.UpdateResult, error) {
	collection := repository.db.GetClient().Database("fitness_db").Collection("routines")

	filter := bson.M{"_id": routine.ID}
	update := bson.M{"$set": bson.M{
		"owner_id":    routine.OwnerID,
		"name":        routine.Name,
		"description": routine.Description,
		"entries":     routine.Entries,
		"is_public":   routine.IsPublic,
		"updated_at":  routine.UpdatedAt,
	}}

	result, err := collection.UpdateOne(context.TODO(), filter, update)
	return result, err
}

func (repository RoutineRepository) DeleteRoutine(id primitive.ObjectID) (*mongo.DeleteResult, error) {
	collection := repository.db.GetClient().Database("fitness_db").Collection("routines")

	filter := bson.M{"_id": id}
	result, err := collection.DeleteOne(context.TODO(), filter)
	return result, err
}
