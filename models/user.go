package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FullName    string             `bson:"full_name" json:"full_name"`
	PhoneNumber string             `bson:"phone_number" json:"phone_number"`
	Email       string             `bson:"email" json:"email"`
	Password    string             `bson:"password" json:"password"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}
