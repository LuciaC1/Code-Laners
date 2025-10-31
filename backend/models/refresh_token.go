package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RefreshToken struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"user_id" json:"userId"`
	Token     string             `bson:"token" json:"token"`
	ExpiresAt time.Time          `bson:"expires_at" json:"expiresAt"`
	Revoked   bool               `bson:"revoked" json:"revoked"`
	CreatedAt time.Time          `bson:"created_at" json:"createdAt"`
	RevokedAt *time.Time         `bson:"revoked_at,omitempty" json:"revokedAt,omitempty"`
}
