package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Book struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Category    string             `bson:"category" json:"category"`
	Description string             `bson:"description" json:"description"`
	Author      string             `bson:"author" json:"author"`
	Quality     uint32             `bson:"quality" json:"quality"`
	Language    string             `bson:"language" json:"language"`
	Price       uint32             `bson:"price" json:"price"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}
