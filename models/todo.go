package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Todo struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Todos string 
	CreatedAt time.Time
	UpdatedAt time.Time
}