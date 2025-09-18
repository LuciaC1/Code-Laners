package repositories

import (
    "context"
	"Code-Laners\backend\models"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
)

type UserRepositoryInterface interface{
	GetUser(nombre string) ([]model.Avion, error)
    GetUserByID(id string) (model.Avion, error)
    CreateUser(avion model.Avion) (*mongo.InsertOneResult, error)
    UpdateUser(avion model.Avion) (*mongo.UpdateResult, error)
    DeleteUser(id primitive.ObjectID) (*mongo.DeleteResult, error)
}
type UserRepository struct{
	db DB
}