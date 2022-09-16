package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Rating struct {
	ID         *primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty""`
	RatedBy    *primitive.ObjectID `json:"rated_by,omitempty" bson:"rated_by,omitempty""`
	ThumbsUp   int                 `json:"thumbsUp" bson:"thumbsUp"`
	ThumbsDown int                 `json:"thumbsDown" bson:"thumbsDown"`
	CreatedAt  time.Time           `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time           `json:"updated_at" bson:"updated_at"`
}
