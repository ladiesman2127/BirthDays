package models

type Session struct {
	AccessToken  *string `bson:"access_token"`
	RefreshToken *string `bson:"refresh_token"`
}

func NewSession(access *string, refresh *string) *Session {
	return &Session{AccessToken: access, RefreshToken: refresh}
}
