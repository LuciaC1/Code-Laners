package repositories

import (
	"context"

	"backend/database"
	"backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepositoryInterface interface {
	GetUser(name string) ([]models.User, error)
	GetUserByID(id string) (models.User, error)
	CreateUser(user models.User) (*mongo.InsertOneResult, error)
	UpdateUser(user models.User) (*mongo.UpdateResult, error)
	DeleteUser(id primitive.ObjectID) (*mongo.DeleteResult, error)
}
type UserRepository struct {
	db database.DB
}

func NewUserRepository(db database.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}
func (repository UserRepository) GetUser(name string) ([]models.User, error) {
	collection := repository.db.GetClient().Database("fitness_db").Collection("users")

	var filter bson.M

	if name != "" {
		filter = bson.M{"name": bson.M{"$regex": name, "$options": "i"}}
	} else {
		filter = bson.M{}
	}

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var users []models.User
	for cursor.Next(context.Background()) {
		var user models.User
		err := cursor.Decode(&user)
		if err != nil {
			continue
		}
		users = append(users, user)
	}

	return users, nil
}
func (repository UserRepository) GetUserByID(id string) (models.User, error) {
	collection := repository.db.GetClient().Database("fitness_db").Collection("users")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.User{}, err
	}

	filter := bson.M{"_id": objectID}
	var user models.User

	err = collection.FindOne(context.TODO(), filter).Decode(&user)
	return user, err
}

func (repository UserRepository) CreateUser(user models.User) (*mongo.InsertOneResult, error) {
	collection := repository.db.GetClient().Database("fitness_db").Collection("users")
	result, err := collection.InsertOne(context.TODO(), user)
	return result, err
}

func (repository UserRepository) UpdaterUser(user models.User) (*mongo.UpdateResult, error) {
	collection := repository.db.GetClient().Database("fitness_db").Collection("users")

	filter := bson.M{"_id": user.ID}
	update := bson.M{"$set": bson.M{
		"name":             user.Name,
		"email":            user.Email,
		"password_hash":    user.PasswordHash,
		"role":             user.Role,
		"date_of_birth":    user.DateOfBirth,
		"weight,omitempty": user.Weight,
		"height,omitempty": user.Height,
		"level,omitempty":  user.Level,
		"goals,omitempty":  user.Goals,
		"created_at":       user.CreatedAt,
		"updated_at":       user.UpdatedAt,
	}}

	result, err := collection.UpdateOne(context.TODO(), filter, update)
	return result, err
}

func (repository UserRepository) DeleteUser(id primitive.ObjectID) (*mongo.DeleteResult, error) {
	collection := repository.db.GetClient().Database("fitness_db").Collection("users")

	filter := bson.M{"_id": id}
	result, err := collection.DeleteOne(context.TODO(), filter)
	return result, err
}
