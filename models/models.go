package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LoginPatient struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Patient struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Email       string             `json:"email" bson:"email"`
	Password    string             `json:"password" bson:"password"`
	Role        string             `json:"role" bson:"role"`
	FullName    string             `json:"full_name" bson:"full_name"`
	DOB         string             `json:"dob" bson:"dob"`
	Phone       string             `json:"phone" bson:"phone"`
	HomeAddress string             `json:"home_address" bson:"home_address"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

type GetPatient struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Email     string             `json:"email" bson:"email"`
	Role      string             `json:"role" bson:"role"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}
