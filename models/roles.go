package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Roles struct {
	ID        *primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty""`
	Name      string              `json:"name" bson:"name"`
	CreatedAt time.Time           `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time           `json:"updated_at" bson:"updated_at"`
}
