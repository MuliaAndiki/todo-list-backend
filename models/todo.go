package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Status string

func (s Status) IsValid() bool {
	switch s {
	case StatusPending, StatusInProgres, StatusDone:
		return true
	}
	return false
}

const (
	StatusPending   Status = "pending"
	StatusInProgres Status = "progress"
	StatusDone      Status = "done"
)

type Todo struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"userId" json:"userId"`
	Todos     string
	Status    Status
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time
}
