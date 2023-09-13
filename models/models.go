package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Book struct {
	ID        uint      `json:"id" bson:"list_id"`
	Title     string    `json:"title" bson:"list_id"`
	Author    string    `json:"author" bson:"list_id"`
	CreatedAt time.Time `json:"created_at" bson:"list_id"`
	UpdatedAt time.Time `json:"updated_at" gormbson:"list_id"`
}

type CreateBook struct {
	Title  string `json:"title" binding:"required"`
	Author string `json:"author" binding:"required"`
}

type UpdateBook struct {
	Title  string `json:"title" bson:"list_id"`
	Author string `json:"author" bson:"list_id"`
}

type LoginUser struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type User struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Email     string             `json:"email" bson:"email"`
	Password  string             `json:"password" bson:"password"`
	Role      string             `json:"role" bson:"role"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}
