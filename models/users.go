package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Users struct {
	ID           *primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty""`
	Name         string              `json:"name"`
	Email        string              `json:"email"`
	Password     string              `json:"password"`
	RefreshToken string              `json:"refresh_token" bson:"refresh_token"`
	UserRole     *primitive.ObjectID `json:"user_role" bson:"user_role"`
	CreatedAt    time.Time           `json:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at"`
}

type UsersResponse struct {
	ID           *primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty""`
	Name         string              `json:"name"`
	Email        string              `json:"email"`
	RefreshToken string              `json:"refresh_token" bson:"refresh_token"`
	UserRole     *primitive.ObjectID `json:"user_role" bson:"user_role"`
	CreatedAt    time.Time           `json:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at"`
}
