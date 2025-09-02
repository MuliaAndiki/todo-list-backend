package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Role string

func (r Role) IsValid() bool {
	switch r {
	case RoleAdmin, RoleUser:
		return true
	}
	return false
}

const (
	RoleAdmin Role = "Admin"
	RoleUser  Role = "User"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Fullname  string             `bson:"fullname" json:"fullname"`
	Email     string             `bson:"email" json:"email"`
	Password  string             `bson:"password,omitempty" json:"-"`
	Role      Role               `bson:"role" json:"role"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}
