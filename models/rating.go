package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Rating struct {
	ID         *primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserId     *primitive.ObjectID `json:"user_id" bson:"user_id,"`
	RatedBy    *primitive.ObjectID `json:"rated_by" bson:"rated_by"`
	ThumbsUp   int64               `json:"thumbs_up" bson:"thumbs_up"`
	ThumbsDown int64               `json:"thumbs_down" bson:"thumbs_down"`
	CreatedAt  time.Time           `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time           `json:"updated_at" bson:"updated_at"`
}
