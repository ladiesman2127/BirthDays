package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID            primitive.ObjectID   `bson:"_id"`
	Phone         *string              `bson:"phone,unique" validate:"required"`
	Login         *string              `bson:"login" validate:"required"`
	Password      *string              `bson:"password" validate:"required"`
	Name          *string              `bson:"name" validate:"required"`
	BirthDay      *string              `bson:"birthday" validate:"required"`
	Subscriptions []primitive.ObjectID `bson:"subscriptions"`
	Subscribers   []primitive.ObjectID `bson:"subscribers"`
	Session       *Session             `bson:"session"`
}
