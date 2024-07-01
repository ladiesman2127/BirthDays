package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ladiesman2127/birthdays/internal/app/models"
)

var SECRET_KEY string = os.Getenv("SECRET_KEY")

func NewSession(login *string) (*models.Session, error) {
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": *login,
		"iss": "birthdayApp",
		"exp": time.Now().Add(time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}).SignedString([]byte(SECRET_KEY))
	if err != nil {
		return nil, err
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": *login,
		"iss": "birthdayApp",
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
		"iat": time.Now().Unix(),
	}).SignedString([]byte(SECRET_KEY))
	if err != nil {
		return nil, err
	}
	return models.NewSession(&accessToken, &refreshToken), nil
}

func ParseToken(tokenToParse *string) (*jwt.Token, error) {
	return jwt.Parse(*tokenToParse, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(SECRET_KEY), nil
	})
}
