package models

type Credentials struct {
	Login    *string `bson:"login" validate:"required"`
	Password *string `bson:"password" validate:"required"`
}

func NewCredentials(login *string, password *string) *Credentials {
	return &Credentials{Login: login, Password: password}
}
