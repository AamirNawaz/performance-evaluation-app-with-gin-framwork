package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Rating struct {
	ID        *primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty""`
	RatedBy   *primitive.ObjectID `json:"rated_by,omitempty" bson:"rated_by,omitempty""`
	ThumsUp   int                 `json:"thumsUp" bson:"thumsUp"`
	ThumsDown int                 `json:"thumsDown" bson:"thumsDown"`
	CreatedAt time.Time           `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time           `json:"updated_at" bson:"updated_at"`
}
