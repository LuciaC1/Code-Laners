package repositories

import (
    "context"
	"Code-Laners/backend/models"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
)

type UserRepositoryInterface interface{
	GetUser(name string) ([]models.User, error)
    GetUserByID(id string) (models.User, error)
    CreateUser(user models.User) (*mongo.InsertOneResult, error)
    UpdateUser(user models.User) (*mongo.UpdateResult, error)
    DeleteUser(id primitive.ObjectID) (*mongo.DeleteResult, error)
}
type UserRepository struct{
	db DB
}
func NewUserRepository(db DB) *UserRepository {
    return &UserRepository{
        db: db,
    }
}
func (repository UserRepository) GetUser(name string)([]models.User, error){
	collection := repository.db.GetClient().Database("fitness_db").collection("users")

	var filter bson.M

	if name != ""{
		filter = bson.M{"name": bson.M{"$regex": name, "$options": "i"}}
	} else {
        filter = bson.M{}
    }
}