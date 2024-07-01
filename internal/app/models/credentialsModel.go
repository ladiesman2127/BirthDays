package models

type Credentials struct {
	Login    *string `bson:"login" validate:"required"`
	Password *string `bson:"password" validate:"required"`
}
