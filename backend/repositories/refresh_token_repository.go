package repositories

import (
	"context"
	"time"

	"backend/database"
	"backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RefreshTokenRepositoryInterface interface {
	Save(token models.RefreshToken) (*mongo.InsertOneResult, error)
	GetByToken(token string) (models.RefreshToken, error)
	Revoke(token string) (*mongo.UpdateResult, error)
	RevokeAllForUser(userID primitive.ObjectID) (*mongo.UpdateResult, error)
}

type RefreshTokenRepository struct {
	db database.DB
}

func NewRefreshTokenRepository(db database.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

func (r RefreshTokenRepository) collection() *mongo.Collection {
	return r.db.GetClient().Database("fitness_db").Collection("refresh_tokens")
}

func (r RefreshTokenRepository) Save(token models.RefreshToken) (*mongo.InsertOneResult, error) {
	if token.ID.IsZero() {
		token.ID = primitive.NewObjectID()
	}
	if token.CreatedAt.IsZero() {
		token.CreatedAt = time.Now()
	}
	return r.collection().InsertOne(context.TODO(), token)
}

func (r RefreshTokenRepository) GetByToken(token string) (models.RefreshToken, error) {
	var rt models.RefreshToken
	filter := bson.M{"token": token}
	err := r.collection().FindOne(context.TODO(), filter).Decode(&rt)
	return rt, err
}

func (r RefreshTokenRepository) Revoke(token string) (*mongo.UpdateResult, error) {
	now := time.Now()
	filter := bson.M{"token": token}
	update := bson.M{"$set": bson.M{"revoked": true, "revoked_at": now}}
	return r.collection().UpdateOne(context.TODO(), filter, update)
}

func (r RefreshTokenRepository) RevokeAllForUser(userID primitive.ObjectID) (*mongo.UpdateResult, error) {
	now := time.Now()
	filter := bson.M{"user_id": userID, "revoked": false}
	update := bson.M{"$set": bson.M{"revoked": true, "revoked_at": now}}

	opts := options.Update().SetUpsert(false)
	res, err := r.collection().UpdateMany(context.TODO(), filter, update, opts)
	if err != nil {
		return nil, err
	}

	return &mongo.UpdateResult{MatchedCount: res.MatchedCount, ModifiedCount: res.ModifiedCount}, nil
}
